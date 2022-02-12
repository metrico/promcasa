package apirouterv1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/metrico/promcasa/config"
	controllerv1 "github.com/metrico/promcasa/controller"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/service"
	"github.com/metrico/promcasa/utils/logger"
	"github.com/patrickmn/go-cache"
)

func RouteInsertDataApis(app fiber.Router, dataSession []*sqlx.DB,
	databaseNodeMap *[]model.DataDatabasesMap, goCache *cache.Cache) {

	// initialize service of user
	insertService := service.InsertService{
		ServiceData:     service.ServiceData{Session: dataSession},
		DatabaseNodeMap: databaseNodeMap,
		GoCache:         goCache,
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

	// initialize user controller
	ict := controllerv1.InsertController{
		InsertService: &insertService,
	}

	//Check DB status
	go doConfigDabaseStats(&insertService)

	//Insert Samples
	go insertService.InsertSamples()
	//Insert Timeseries
	go insertService.InsertTimeSeries()

	insertService.ReloadFingerprints()

	// push streams
	app.Post("/push", ict.PushStream)

}

// make a ping keep alive
func doConfigDabaseStats(us *service.InsertService) {

	for {
		logger.Debug("Starting config database check")
		us.RunWatcherConfigDatabaseStats()
		time.Sleep(time.Duration(60) * time.Second)
	}
}
