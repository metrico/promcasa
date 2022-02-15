package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/Jeffail/gabs/v2"
	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/utils/async"
	"github.com/metrico/promcasa/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type InsertService struct {
	ServiceData
	DatabaseNodeMap *[]model.DataDatabasesMap
}

// internal sync
func (ss *InsertService) RunWatcherConfigDatabaseStats() error {

	//var searchData
	for idx, db := range ss.Session {
		logger.Debug("RunWatcherConfigDatabaseStats: CHECK DataDB: ", (*ss.DatabaseNodeMap)[idx].Name)

		if err := db.Ping(); err != nil {
			(*ss.DatabaseNodeMap)[idx].Online = false
			logger.Debug("node is offline: ", (*ss.DatabaseNodeMap)[idx].Name)
			if exception, ok := err.(*clickhouse.Exception); ok {
				logger.Error(fmt.Sprintf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
			} else {
				logger.Debug("ping db data ", err)
			}
		} else {
			logger.Debug("node is online: ", (*ss.DatabaseNodeMap)[idx].Name)
			(*ss.DatabaseNodeMap)[idx].Online = true

		}
	}

	return nil
}

// this method create new user in the database
// it doesn't check internally whether all the validation are applied or not
func (ss *InsertService) DoMetricsQueries() error {

	if !(*ss.DatabaseNodeMap)[config.Setting.CurrentDataBaseIndex].Online {
		logger.Error("the node is offline:")
		return fmt.Errorf("the node is offline")
	}

	currentDateTime := time.Now()

	var ResultFuture []async.Future

	for index, qObj := range config.Setting.DATABASE_METRICS {

		dbIndex := config.Setting.CurrentDataBaseIndex

		if qObj.MetricLiveView {
			logger.Debug("This is live view, skip it")
			continue
		}

		if qObj.LastTime.Add(qObj.RefreshTimeout).After(currentDateTime) {
			logger.Debug("Not yet.....")
			continue
		}

		sqlQuery := strings.Replace(qObj.Query, "$refresh", fmt.Sprintf("interval %d SECOND", int(qObj.RefreshTimeout.Seconds())), -1)
		logger.Debug("Execute query: ", index, sqlQuery)

		future := async.ExecAsyncSql(func(query string, lIndex uint, qIndex int) model.AsyncSqlResult {
			logger.Debug("Execute Async process on node: ", lIndex, qIndex)
			result := model.AsyncSqlResult{QueryIndex: qIndex}
			result.Rows, result.Err = ss.Session[lIndex].Queryx(query) // (*sql.Rows, error)
			return result
		}, sqlQuery, dbIndex, index)

		config.Setting.DATABASE_METRICS[index].LastTime = currentDateTime

		ResultFuture = append(ResultFuture, future)
	}

	for index := range ResultFuture {

		result := ResultFuture[index].Await()
		var objects []map[string]interface{}

		if result.Err != nil {
			logger.Error("Error in future process:", result.Err)
			logger.Error("Error at index: ", index)
			continue
		}

		rows := result.Rows
		defer rows.Close()

		for rows.Next() {
			// figure out what columns were returned
			// the column names will be the JSON object field keys
			columns, err := rows.ColumnTypes()
			if err != nil {
				logger.Error("bad column types: ", err.Error())
				return err
			}

			// Scan needs an array of pointers to the values it is setting
			// This creates the object and sets the values correctly
			values := make([]interface{}, len(columns))
			object := map[string]interface{}{}
			for i, column := range columns {
				v := reflect.New(column.ScanType()).Interface()
				switch v.(type) {
				case *[]uint8:
					v = new(string)
				default:
					// use this to find the type for the field
					// you need to change
					//logger.Debug("%v: %T", column.Name(), v)
				}

				nn := strings.Replace(column.Name(), ".", "->", -1)
				object[nn] = v
				values[i] = object[nn]
			}

			err = rows.Scan(values...)
			if err != nil {
				logger.Error("found error during scan of object:", err.Error())
			} else {
				objects = append(objects, object)
			}
		}

		rowsObj, _ := json.Marshal(objects)
		data, _ := gabs.ParseJSON(rowsObj)

		//go is static type language - in this case for each type we have to do own block

		promCasaMetric := config.Setting.DATABASE_METRICS[result.QueryIndex]
		var ok bool
		if promCasaMetric.MetricType == "gauge" {
			var gaugeVec *prometheus.GaugeVec
			if gaugeVec, ok = config.Setting.PromGaugeMap[promCasaMetric.Name]; ok {
				//do something here
			} else {

				if config.Setting.PromGaugeMap == nil {
					config.Setting.PromGaugeMap = make(map[string]*prometheus.GaugeVec)
				}

				config.Setting.PromGaugeMap[promCasaMetric.Name] = promauto.NewGaugeVec(prometheus.GaugeOpts{
					Name: promCasaMetric.Name,
					Help: promCasaMetric.Help},
					promCasaMetric.MetricLabels)

				gaugeVec = config.Setting.PromGaugeMap[promCasaMetric.Name]

			}

			for _, value := range data.Children() {

				var counter float64
				labels := prometheus.Labels{}

				for _, label := range promCasaMetric.MetricLabels {

					if value.Exists(label) {
						labels[label] = value.S(label).Data().(string)
					}
				}

				if value.Exists(promCasaMetric.CounterName) && value.S(promCasaMetric.CounterName).Data() != nil {
					counter = value.S(promCasaMetric.CounterName).Data().(float64)
				}

				gaugeVec.With(labels).Set(counter)
			}
		} else if promCasaMetric.MetricType == "histogram" {
			var histogramVec *prometheus.HistogramVec
			if histogramVec, ok = config.Setting.PromHistogramMap[promCasaMetric.Name]; ok {
				//do something here
			} else {

				if config.Setting.PromHistogramMap == nil {
					config.Setting.PromHistogramMap = make(map[string]*prometheus.HistogramVec)
				}

				config.Setting.PromHistogramMap[promCasaMetric.Name] = promauto.NewHistogramVec(prometheus.HistogramOpts{
					Name:    promCasaMetric.Name,
					Help:    promCasaMetric.Help,
					Buckets: []float64{0.01, 0.02, 0.05, 0.1},
				},
					promCasaMetric.MetricLabels)

				histogramVec = config.Setting.PromHistogramMap[promCasaMetric.Name]

			}

			for _, value := range data.Children() {

				var counter float64
				labels := prometheus.Labels{}

				for _, label := range promCasaMetric.MetricLabels {

					if value.Exists(label) {
						labels[label] = value.S(label).Data().(string)
					}
				}

				if value.Exists(promCasaMetric.CounterName) && value.S(promCasaMetric.CounterName).Data() != nil {
					counter = value.S(promCasaMetric.CounterName).Data().(float64)
				}

				histogramVec.With(labels).Observe(counter)
			}
		} else if promCasaMetric.MetricType == "counter" {
			var counterVec *prometheus.CounterVec
			if counterVec, ok = config.Setting.PromCounterMap[promCasaMetric.Name]; ok {
				//do something here
			} else {

				if config.Setting.PromCounterMap == nil {
					config.Setting.PromCounterMap = make(map[string]*prometheus.CounterVec)
				}

				config.Setting.PromCounterMap[promCasaMetric.Name] = promauto.NewCounterVec(prometheus.CounterOpts{
					Name: promCasaMetric.Name,
					Help: promCasaMetric.Help,
				},
					promCasaMetric.MetricLabels)

				counterVec = config.Setting.PromCounterMap[promCasaMetric.Name]

			}

			for _, value := range data.Children() {

				var counter float64
				labels := prometheus.Labels{}

				for _, label := range promCasaMetric.MetricLabels {

					if value.Exists(label) {
						labels[label] = value.S(label).Data().(string)
					}
				}

				if value.Exists(promCasaMetric.CounterName) && value.S(promCasaMetric.CounterName).Data() != nil {
					counter = value.S(promCasaMetric.CounterName).Data().(float64)
				}
				counterVec.With(labels).Add(counter)
			}
		}

	}

	return nil
}
