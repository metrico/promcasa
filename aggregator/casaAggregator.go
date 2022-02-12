package aggregator

import (
	"time"

	"github.com/jmoiron/sqlx"
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

	//Check DB status
	go doConfigDabaseStats(&insertService)

	//Check query status
	doMetricsScheduler(&insertService)
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

	for {
		logger.Debug("Starting queries check")
		us.DoMetricsQueries()
		time.Sleep(time.Duration(60) * time.Second)
	}
}
