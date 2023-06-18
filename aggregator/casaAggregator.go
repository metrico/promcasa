package aggregator

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/service"
	"github.com/metrico/promcasa/utils/logger"
)

func ActivateTimer(dataSession []*sqlx.DB, databaseNodeMap *[]model.DataDatabasesMap) {

	// initialize service of user
	insertService := service.InsertService{
		ServiceData:     service.ServiceData{Session: dataSession},
		DatabaseNodeMap: databaseNodeMap,
	}

	//Check DB status - DB watcher is true
	if config.Setting.SYSTEM_SETTINGS.DBWatcher {
		go doConfigDabaseStats(&insertService)
	}

	//Check query status
	go doMetricsScheduler(&insertService)
}

// make a ping keep alive
func doConfigDabaseStats(us *service.InsertService) {

	for {
		logger.Debug("Starting config database check")
		us.RunWatcherConfigDatabaseStats()
		time.Sleep(time.Duration(60) * time.Second)
	}
}

// make a ping keep alive
func doMetricsScheduler(us *service.InsertService) {

	refreshTimeout, _ := time.ParseDuration(config.Setting.SYSTEM_SETTINGS.SystemRefreshCheck)

	for {
		logger.Debug("Starting queries check")
		us.DoMetricsQueries()
		time.Sleep(refreshTimeout)
	}
}
