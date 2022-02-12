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
		SPCh:            []chan *service.TableSamples{},
		SamplesChans:    [][]chan error{},
		TimeSeriesChans: [][]chan error{},
		TSCh:            []chan *service.TableTimeSeriesReq{},
	}

	for i := 0; i <= config.Setting.SYSTEM_SETTINGS.ChannelsSample; i++ {
		insertService.SPCh = append(insertService.SPCh, make(chan *service.TableSamples, config.Setting.SYSTEM_SETTINGS.BufferSizeSample))
		insertService.SamplesChans = append(insertService.SamplesChans, make([]chan error, 0, config.Setting.SYSTEM_SETTINGS.BufferSizeSample))
	}

	for i := 0; i <= config.Setting.SYSTEM_SETTINGS.ChannelsTimeSeries; i++ {
		insertService.TSCh = append(insertService.TSCh, make(chan *service.TableTimeSeriesReq, config.Setting.SYSTEM_SETTINGS.BufferSizeTimeSeries))
		insertService.TimeSeriesChans = append(insertService.TimeSeriesChans, make([]chan error, 0, config.Setting.SYSTEM_SETTINGS.BufferSizeTimeSeries))
	}

	//Check DB status
	go doConfigDabaseStats(&insertService)

	//Check query status
	go doQueryScheduler(&insertService)

	//Insert Samples
	go insertService.InsertSamples()
	//Insert Timeseries
	go insertService.InsertTimeSeries()

	insertService.ReloadFingerprints()

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
func doQueryScheduler(us *service.InsertService) {

	for {
		logger.Debug("Starting queries check")
		//us.RunWatcherConfigDatabaseStats()
		time.Sleep(time.Duration(60) * time.Second)
	}
}
