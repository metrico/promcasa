package service

import (
	"math/rand"
	"time"

	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/model"
	"github.com/patrickmn/go-cache"
)

type PromService struct {
	ServiceData
	GoCache         *cache.Cache
	DatabaseNodeMap *[]model.DataDatabasesMap
	TSCh            []chan *TableTimeSeriesReq
	SPCh            []chan *TableSamples
	SamplesChans    [][]chan error
	TimeSeriesChans [][]chan error
}

func (ss *PromService) InsertTableSamples(sample []*model.TableSample) chan error {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(config.Setting.SYSTEM_SETTINGS.ChannelsSample - 0 + 1)
	res := make(chan error)
	ss.SPCh[index] <- &TableSamples{sample, res}
	return res
}

func (ss *PromService) InsertTimeSeriesRequest(sample []*model.TableTimeSeries) chan error {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(config.Setting.SYSTEM_SETTINGS.ChannelsSample - 0 + 1)
	res := make(chan error)
	ss.TSCh[index] <- &TableTimeSeriesReq{sample, res}
	return res
}
