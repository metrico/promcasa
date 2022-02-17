package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/metrico/promcasa/aggregator"
	"github.com/metrico/promcasa/utils/helpers"
	"github.com/metrico/promcasa/utils/jobqueue"
	"github.com/metrico/promcasa/utils/promcasautils"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mcuadros/go-defaults"
	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/utils/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

//params for  Services
type ServicesObject struct {
	dataDBSession   []*sqlx.DB
	databaseNodeMap []model.DataDatabasesMap
}

type CustomValidator struct {
	validator *validator.Validate
}

// validate function
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var appFlags CommandLineFlags
var servicesObject ServicesObject

//params for Flags
type CommandLineFlags struct {
	InitializeDB    *bool   `json:"initialize_db"`
	ShowHelpMessage *bool   `json:"help"`
	ShowVersion     *bool   `json:"version"`
	ConfigPath      *string `json:"config_path"`
	LogPath         *string `json:"log_path"`
	LogName         *string `json:"log_name"`
}

/* init flags */
func initFlags() {
	appFlags.InitializeDB = flag.Bool("initialize_db", false, "initialize the database and create all tables")
	appFlags.ShowHelpMessage = flag.Bool("help", false, "show help")
	appFlags.ShowVersion = flag.Bool("version", false, "show version")
	appFlags.ConfigPath = flag.String("config-path", "/usr/local/promcasa/etc", "the path to the promcasaapp config file")
	appFlags.LogName = flag.String("log-name", "", "the name prefix of the log file.")
	appFlags.LogPath = flag.String("log-path", "", "the path for the log file.")

	flag.Parse()
}

func main() {

	//init flags
	initFlags()

	cfg := new(config.PromCasaSettingServer)
	defaults.SetDefaults(cfg) //<-- This set the defaults values
	config.Setting = *cfg

	/* first check admin flags */
	checkHelpVersionFlags()

	promcasautils.Colorize(promcasautils.ColorRed, promcasautils.PromCasaLogo)

	//ReadConfig
	readConfig()

	//SystemParams
	setFastConfigSettings()

	//Set to max cpu if the value is equals 0
	if config.Setting.SYSTEM_SETTINGS.CPUMaxProcs == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(config.Setting.SYSTEM_SETTINGS.CPUMaxProcs)
	}

	// initialize logger
	logger.InitLogger()

	// configure new db session
	servicesObject.dataDBSession, servicesObject.databaseNodeMap = getDataDBSession()
	for val := range servicesObject.dataDBSession {
		defer servicesObject.dataDBSession[val].Close()
	}

	if len(servicesObject.dataDBSession) == 0 {
		promcasautils.Colorize(promcasautils.ColorRed, "\r\nWe don't have any active DB session configured. Please check your config\r\n")
		os.Exit(0)
	}

	/* init job push queue channel */
	if config.Setting.PROMETHEUS_CLIENT.EnablePush {
		jobqueue.InitQueue(config.Setting.PROMETHEUS_CLIENT.QueueJobELements)
	}

	//Api
	// configure to serve WebServices
	configureAsHTTPServer()
}

func checkHelpVersionFlags() {
	if *appFlags.ShowHelpMessage {
		flag.Usage()
		os.Exit(0)
	}

	if *appFlags.ShowVersion {
		fmt.Printf("VERSION: %s\r\n", VERSION_APPLICATION)
		os.Exit(0)
	}
}

//https://github.com/atreugo/examples/blob/master/basic/main.go
//https://github.com/jackwhelpton/fasthttp-routing
func readConfig() {
	// Getting constant values
	if configEnv := os.Getenv("PROMCASA_APPENV"); configEnv != "" {
		viper.SetConfigName("promcasa_" + configEnv)
	} else {
		viper.SetConfigName("promcasa")
	}
	viper.SetConfigType("json")

	if configPath := os.Getenv("PROMCASA_APPPATH"); configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(*appFlags.ConfigPath)
	}

	viper.AddConfigPath(".")

	//Default value
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - checking env: ", err)
		logger.Error("No configuration file loaded - using defaults - checking env")
	}

	viper.SetConfigName("promcasa_custom")
	err = viper.MergeInConfig()
	if err != nil {
		logger.Debug("No custom configuration file loaded.")
	}

	//Env variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "[", "_", "]", ""))
	viper.SetEnvPrefix(config.Setting.EnvPrefix)
	SetEnvironDataBase()

	//Bind Env from Config
	BindEnvs(config.Setting)

	err = viper.Unmarshal(&config.Setting, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		logger.Debug("couldn't unmarshal viper.")
	}

	/* database data connection */
	/* database queries */

	var reData = regexp.MustCompile(`database_data\[(\d)\]`)
	var reQueries = regexp.MustCompile(`database_metrics\[(\d)\]`)

	envParamsData := []int{}
	envParamsQueries := []int{}
	allSettings := viper.AllSettings()
	for key, _ := range allSettings {
		if strings.HasPrefix(key, "database_data[") {
			key = reData.ReplaceAllString(key, "$1")
			i, err := strconv.Atoi(key)
			if err == nil {
				envParamsData = append(envParamsData, i)
			}
		} else if strings.HasPrefix(key, "database_queries[") {
			key = reQueries.ReplaceAllString(key, "$1")
			i, err := strconv.Atoi(key)
			if err == nil {
				envParamsQueries = append(envParamsQueries, i)
			}
		}
	}

	//Read the data db connection from config. This is fix because defaults doesn't know the size of array
	if viper.IsSet("database_data") {
		config.Setting.DATABASE_DATA = nil
		dataConfig := viper.Get("database_data")
		dataVal := dataConfig.([]interface{})
		for idx := range dataVal {
			val := dataVal[idx].(map[string]interface{})
			data := config.PromCasaDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_DATA = append(config.Setting.DATABASE_DATA, data)
		}
	}

	//We should do extraction and after sorting 0,1,2,3
	sort.Ints(envParamsData[:])
	//Here we do ENV check
	for _, idx := range envParamsData {
		value := allSettings[fmt.Sprintf("database_data[%d]", idx)]
		val := value.(map[string]interface{})
		//If the configuration already exists - we replace only existing params
		if len(config.Setting.DATABASE_DATA) > idx {
			err := mapstructure.Decode(val, &config.Setting.DATABASE_DATA[idx])
			if err != nil {
				logger.Error("ERROR during mapstructure decode[0]:", err)
			}
		} else {
			data := config.PromCasaDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_DATA = append(config.Setting.DATABASE_DATA, data)
		}
	}

	//queries
	if viper.IsSet("database_metrics") {
		config.Setting.DATABASE_METRICS = nil
		dataConfig := viper.Get("database_metrics")
		dataVal := dataConfig.([]interface{})
		for idx := range dataVal {
			val := dataVal[idx].(map[string]interface{})
			data := config.PromCasaMetrics{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_METRICS = append(config.Setting.DATABASE_METRICS, data)
		}
	}

	//We should do extraction and after sorting 0,1,2,3
	sort.Ints(envParamsQueries[:])
	//Here we do ENV check
	for _, idx := range envParamsQueries {
		value := allSettings[fmt.Sprintf("database_metrics[%d]", idx)]
		val := value.(map[string]interface{})
		//If the configuration already exists - we replace only existing params
		if len(config.Setting.DATABASE_METRICS) > idx {
			err := mapstructure.Decode(val, &config.Setting.DATABASE_METRICS[idx])
			if err != nil {
				logger.Error("ERROR during mapstructure decode[0]:", err)
			}
		} else {
			data := config.PromCasaMetrics{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_METRICS = append(config.Setting.DATABASE_METRICS, data)
		}
	}

	//Set to 0
	for index, qObj := range config.Setting.DATABASE_METRICS {
		config.Setting.DATABASE_METRICS[index].RefreshTimeout, _ = time.ParseDuration(qObj.RefreshString)
	}

	//Check the command line
	if *appFlags.LogName != "" {
		config.Setting.LOG_SETTINGS.Name = *appFlags.LogName
	}

	if *appFlags.LogPath != "" {
		config.Setting.LOG_SETTINGS.Path = *appFlags.LogPath
	}

	//viper.Debug()
}

//system params for replications, groups
func setFastConfigSettings() {

	/***********************************/

	minVersion := config.Setting.HTTPS_SETTINGS.MinTLSVersionString

	if minVersion == "TLS1.0" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS10
	} else if minVersion == "TLS1.1" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS11
	} else if minVersion == "TLS1.2" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS12
	} else if minVersion == "TLS1.3" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS13
	}

	maxVersion := config.Setting.HTTPS_SETTINGS.MaxTLSVersionString

	if maxVersion == "TLS1.0" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS10
	} else if maxVersion == "TLS1.1" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS11
	} else if maxVersion == "TLS1.2" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS12
	} else if maxVersion == "TLS1.3" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS13
	}
}

// getSession creates a new mongo session and panics if connection error occurs
func getDataDBSession() ([]*sqlx.DB, []model.DataDatabasesMap) {

	dbMap := []*sqlx.DB{}
	dbNodeMap := []model.DataDatabasesMap{}

	// Rlogs
	if logger.RLogs != nil {
		clickhouse.SetLogOutput(logger.RLogs)
	}

	for _, dbObject := range config.Setting.DATABASE_DATA {

		timeReadTimeout, _ := time.ParseDuration(dbObject.ReadTimeout)
		timeWriteTimeout, _ := time.ParseDuration(dbObject.WriteTimeout)

		logger.Info(fmt.Sprintf("Connecting to Host: [%s], User:[%s], Name:[%s], Node:[%s], Port:[%d], Timeout: [%s, %s]\n",
			dbObject.Host, dbObject.User, dbObject.Name, dbObject.Node,
			dbObject.Port, dbObject.ReadTimeout, dbObject.WriteTimeout))

		connectString := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&read_timeout=%d&write_timeout=%d&compress=true&debug=%t",
			dbObject.Host, dbObject.Port, dbObject.User, dbObject.Password, dbObject.Name,
			int(timeReadTimeout.Seconds()), int(timeWriteTimeout.Seconds()), dbObject.Debug)

		db, err := sqlx.Open("clickhouse", connectString)
		dbMap = append(dbMap, db)
		dbNodeMap = append(dbNodeMap, model.DataDatabasesMap{Name: dbObject.Name, DBname: dbObject.Name, Host: dbObject.Host,
			TableSeries: dbObject.TableSeries, TableSamples: dbObject.TableSamples, Online: false})
		if err != nil {
			logger.Error(fmt.Sprintf("couldn't make connection to [Host: %s, Node: %s, Port: %d]: \n", dbObject.Host, dbObject.Node, dbObject.Port), err)
			continue
		}

		//Set some limit
		db.SetMaxIdleConns(dbObject.MaxIdleConn)
		db.SetMaxOpenConns(dbObject.MaxOpenConn)

		logger.Info("----------------------------------- ")
		logger.Info("*** Database Session created *** ")
		logger.Info("----------------------------------- ")
	}

	return dbMap, dbNodeMap
}

func configureAsHTTPServer() {

	helpers.SetGlobalLimit(50 * 1024 * 1024)

	httpURL := fmt.Sprintf("%s:%d", config.Setting.HTTP_SETTINGS.Host, config.Setting.HTTP_SETTINGS.Port)
	configFiber := fiber.Config{
		Prefork:           config.Setting.HTTP_SETTINGS.Prefork,
		StreamRequestBody: true,
		//Addr:    httpURL,
	}

	config.Setting.Validate = validator.New()

	serverFiber := fiber.New(configFiber)

	/* check if we should enable prometheus client here */
	if config.Setting.PROMETHEUS_CLIENT.Enable {
		prometheus := fiberprometheus.New(config.Setting.PROMETHEUS_CLIENT.ServiceName)
		prometheus.RegisterAt(serverFiber, config.Setting.PROMETHEUS_CLIENT.MetricsPath)
		serverFiber.Use(prometheus.Middleware)
	}

	/* fractionLost := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "custom_data_fraction_lost",
		Help: "fraction lost"},
		[]string{"node_id"})
	fractionLost.WithLabelValues("test").Inc()
	*/

	runSchedularPopulation()

	if err := serverFiber.Listen(httpURL); err != nil {
		panic(err)
	}

}

func runSchedularPopulation() {

	aggregator.ActivateTimer(servicesObject.dataDBSession, &servicesObject.databaseNodeMap)
}

//this function will check PROMCASA_DATABASE_DATA and set internal bind for viper
//i.e. PROMCASA_DATABASE_DATA_0_HOSTNAME -> database_data[0].hostname
func SetEnvironDataBase() bool {
	var re = regexp.MustCompile(`_(\d)_`)
	for _, s := range os.Environ() {
		if strings.HasPrefix(s, config.Setting.EnvPrefix+"_DATABASE_DATA") {
			a := strings.Split(s, "=")
			key := strings.TrimPrefix(a[0], config.Setting.EnvPrefix+"_")
			key = re.ReplaceAllString(key, "[$1].")
			viper.BindEnv(key)
		}
	}
	return true
}

//Now we should bind the ENV params
func BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		case reflect.Slice:
			continue
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
