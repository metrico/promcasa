package config

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/go-playground/validator.v9"
)

var Setting PromCasaSettingServer

// NAME
var NAME_APPLICATION = "promcasa"

type PromCasaDataBase struct {
	Proto        string `json:"proto" mapstructure:"proto" default:"tcp"`
	User         string `json:"user" mapstructure:"user" default:"promcasa_user"`
	Node         string `json:"node" mapstructure:"node" default:"promcasanode"`
	Password     string `json:"pass" mapstructure:"pass" default:"promcasa_pass"`
	Name         string `json:"name" mapstructure:"name" default:"promcasa_data"`
	Host         string `json:"host" mapstructure:"host" default:"127.0.0.1"`
	TableSamples string `json:"table_samples" mapstructure:"table_samples" default:"samples_v2"`
	TableSeries  string `json:"table_series" mapstructure:"table_series" default:"time_series"`
	Debug        bool   `json:"debug" mapstructure:"debug" default:"false"`
	Port         uint32 `json:"port" mapstructure:"port" default:"9000"`
	ReadTimeout  string `json:"read_timeout" mapstructure:"read_timeout" default:"5m"`
	DialTimeout  string `json:"dial_timeout" mapstructure:"dial_timeout" default:"30s"`
	MaxIdleConn  int    `json:"max_idle_connection" mapstructure:"max_idle_connection" default:"5"`
	MaxOpenConn  int    `json:"max_open_connection" mapstructure:"max_open_connection" default:"50"`
	Primary      bool   `json:"primary" mapstructure:"primary" default:"false"`
	Strategy     string `json:"strategy" mapstructure:"strategy" default:"failover"`
}

type PromCasaMetrics struct {
	Name           string        `json:"name" mapstructure:"name" default:""`
	Help           string        `json:"help" mapstructure:"help" default:""`
	Query          string        `json:"query" mapstructure:"query" default:""`
	CounterName    string        `json:"counter_name" mapstructure:"counter_name" default:"counter"`
	RefreshString  string        `json:"refresh" mapstructure:"refresh" default:"60s"`
	MetricType     string        `json:"type" mapstructure:"type" default:"gauge"`
	MetricLiveView bool          `json:"live_view" mapstructure:"live_view" default:"false"`
	MetricLabels   []string      `json:"labels" mapstructure:"labels" default:"gauge"`
	LastTime       time.Time     `json:"last_time" default:""`
	RefreshTimeout time.Duration `json:"refresh_timeout" default:""`
}

type PromCasaSettingServer struct {
	SrartTime                time.Time `default:""`
	DataBaseStrategy         uint      `default:"0"`
	CurrentDataBaseIndex     uint      `default:"0"`
	DataDatabaseGroupNodeMap map[string][]string
	Validate                 *validator.Validate
	EnvPrefix                string `default:"PROMCASA"`
	PromGaugeMap             map[string]*prometheus.GaugeVec
	PromHistogramMap         map[string]*prometheus.HistogramVec
	PromCounterMap           map[string]*prometheus.CounterVec

	DATABASE_DATA    []PromCasaDataBase `json:"database_data" mapstructure:"database_data"`
	DATABASE_METRICS []PromCasaMetrics  `json:"database_metrics" mapstructure:"database_metrics"`

	SYSTEM_SETTINGS struct {
		HostName             string `json:"hostname" mapstructure:"hostname" default:"hostname"`
		EnableUserAuditLogin bool   `json:"user_audit_login" mapstructure:"user_audit_login" default:"true"`
		UUID                 string `json:"uuid" mapstructure:"uuid" default:""`
		DBBulk               int    `json:"db_bulk" mapstructure:"db_bulk" default:"40000"`
		DBTimer              string `json:"db_timer" mapstructure:"db_timer" default:"1s"`
		DBWatcher            bool   `json:"db_watcher" mapstructure:"db_watcher" default:"true"`
		BufferSizeSample     uint32 `json:"buffer_size_sample" mapstructure:"buffer_size_sample" default:"200000"`
		BufferSizeTimeSeries uint32 `json:"buffer_size_timeseries" mapstructure:"buffer_size_timeseries" default:"200000"`
		CPUMaxProcs          int    `json:"cpu_max_procs" mapstructure:"cpu_max_procs" default:"1"`
		SystemRefreshCheck   string `json:"system_refresh" mapstructure:"system_refresh" default:"10s"`
	} `json:"system_settings" mapstructure:"system_settings"`

	AUTH_SETTINGS struct {
		AuthTokenHeader string `json:"token_header" mapstructure:"token_header" default:"Auth-Token"`
		AuthTokenExpire string `json:"token_expire" mapstructure:"token_expire" default:"1200s"`
	} `json:"auth_settings" mapstructure:"auth_settings"`

	API_SETTINGS struct {
		EnableForceSync   bool `json:"sync_map_force" mapstructure:"sync_map_force" default:"false"`
		EnableTokenAccess bool `json:"enable_token_access" mapstructure:"enable_token_access" default:"true"`
	} `json:"api_settings" mapstructure:"api_settings"`

	HTTP_SETTINGS struct {
		Host       string `json:"host" mapstructure:"host" default:"0.0.0.0"`
		Port       int    `json:"port" mapstructure:"port" default:"3200"`
		Prefork    bool   `json:"prefork" mapstructure:"prefork" default:"false"`
		Gzip       bool   `json:"gzip" mapstructure:"gzip" default:"true"`
		GzipStatic bool   `json:"gzip_static" mapstructure:"gzip_static" default:"true"`
		Debug      bool   `json:"debug" mapstructure:"debug" default:"false"`
		WebSocket  struct {
			Enable bool `json:"enable" mapstructure:"enable" default:"false"`
		} `json:"websocket" mapstructure:"websocket"`
		Enable bool `json:"enable" mapstructure:"enable" default:"true"`
	} `json:"http_settings" mapstructure:"http_settings"`

	HTTPS_SETTINGS struct {
		Host                string `json:"host" mapstructure:"host" default:"0.0.0.0"`
		Port                int    `json:"port" mapstructure:"port" default:"3201"`
		Cert                string `json:"cert" mapstructure:"cert" default:""`
		Key                 string `json:"key" mapstructure:"key" default:""`
		HttpRedirect        bool   `json:"http_redirect" mapstructure:"http_redirect" default:"false"`
		Enable              bool   `json:"enable" mapstructure:"enable" default:"false"`
		MinTLSVersionString string `json:"min_tls_version" mapstructure:"min_tls_version" default:"0"`
		MaxTLSVersionString string `json:"max_tls_version" mapstructure:"max_tls_version" default:"0"`
		MinTLSVersion       uint16 `default:"0"`
		MaxTLSVersion       uint16 `default:"0"`
	} `json:"https_settings" mapstructure:"https_settings"`

	LOG_SETTINGS struct {
		Enable        bool   `json:"enable" mapstructure:"enable" default:"true"`
		MaxAgeDays    uint32 `json:"max_age_days" mapstructure:"max_age_days" default:"7"`
		RotationHours uint32 `json:"rotation_hours" mapstructure:"rotation_hours" default:"24"`
		Path          string `json:"path" mapstructure:"path" default:"./"`
		Level         string `json:"level" mapstructure:"level" default:"error"`
		Name          string `json:"name" mapstructure:"name" default:"promcasa.log"`
		Stdout        bool   `json:"stdout" mapstructure:"stdout" default:"false"`
		Json          bool   `json:"json" mapstructure:"json" default:"true"`
		SysLogLevel   string `json:"syslog_level" mapstructure:"syslog_level" default:"LOG_INFO"`
		SysLog        bool   `json:"syslog" mapstructure:"syslog" default:"false"`
		SyslogUri     string `json:"syslog_uri" mapstructure:"syslog_uri" default:""`
	} `json:"log_settings" mapstructure:"log_settings"`

	PROMETHEUS_CLIENT struct {
		PushURL          string   `json:"push_url" mapstructure:"push_url" default:""`
		TargetIP         string   `json:"target_ip" mapstructure:"target_ip" default:""`
		EnablePush       bool     `json:"push_enable" mapstructure:"push_enable" default:"true"`
		PushInterval     string   `json:"push_interval" mapstructure:"push_interval" default:"60s"`
		PushName         string   `json:"push_name" mapstructure:"push_name" default:""`
		ServiceName      string   `json:"service_name" mapstructure:"service_name" default:"prometheus"`
		MetricsPath      string   `json:"metrics_path" mapstructure:"metrics_path" default:"/metrics"`
		Enable           bool     `json:"enable" mapstructure:"enable" default:"true"`
		AllowIP          []string `json:"allow_ip" mapstructure:"allow_ip" default:"[127.0.0.1]"`
		QueueJobELements uint32   `json:"queue_size" mapstructure:"queue_size" default:"5000"`
	} `json:"prometheus_client" mapstructure:"prometheus_client"`
}
