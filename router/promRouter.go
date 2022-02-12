package apirouterv1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	controllerv1 "github.com/metrico/promcasa/controller"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/service"
	"github.com/patrickmn/go-cache"
)

func RoutePromDataApis(app fiber.Router, dataSession []*sqlx.DB, databaseNodeMap *[]model.DataDatabasesMap, goCache *cache.Cache) {

	// initialize service of user
	promService := service.PromService{
		ServiceData:     service.ServiceData{Session: dataSession},
		DatabaseNodeMap: databaseNodeMap,
		GoCache:         goCache,
		SPCh:            []chan *service.TableSamples{},
		SamplesChans:    [][]chan error{},
		TimeSeriesChans: [][]chan error{},
		TSCh:            []chan *service.TableTimeSeriesReq{},
	}

	// initialize user controller
	ict := controllerv1.PromController{
		PromService: &promService,
	}

	// write prometheus streams
	app.Post("v1/prom/remote/write", ict.WriteStream)

	// Alias from write
	app.Post("prom/remote/write", ict.WriteStream)
}
