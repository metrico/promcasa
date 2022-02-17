package jobqueue

import (
	"fmt"
	"sync"
	"time"

	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type MetricType int

var prometheusPusher *push.Pusher

const (
	GaugeType MetricType = iota
	CounterType
	HistogramType
)

// swagger:model OLDAliasStruct
type JobObject struct {
	NameStats  string
	Type       MetricType
	CreateDate time.Time
}

var (
	// wg is used to force the application wait for all goroutines to finish before exiting.
	queueGroup sync.WaitGroup
	// jobChan is a buffered channel that has the capacity of maximum "elements" resource slot.
	jobChannel chan JobObject
)

// queueJob puts job into channel. If channel buffer is full, return false.
func InitQueue(elements uint32) bool {

	logger.Debug("Starting a job queue with elements: " + fmt.Sprint(elements))
	jobChannel = make(chan JobObject, elements)
	queueGroup.Add(1)

	//Create a pusher
	prometheusPusher = push.New(config.Setting.PROMETHEUS_CLIENT.PushURL, config.Setting.PROMETHEUS_CLIENT.PushName)

	// Run 1 worker to handle jobs.
	go worker(jobChannel, &queueGroup)

	logger.Debug("Started the job queue")

	return true
}

// queueJob puts job into channel. If channel buffer is full, return false.
func CloseQueue() bool {

	logger.Debug("Stopping the job queue")

	close(jobChannel)

	// Block exiting until all the goroutines are finished.
	queueGroup.Wait()

	logger.Debug("Finished the job queue")

	return true
}

// send a job worker processes jobs.
func SendJob(metricName string, metricTytpe MetricType) bool {

	/* copy */
	newJob := JobObject{NameStats: metricName, Type: metricTytpe}
	logger.Debug("Worker pushing job: ", newJob)

	if !queueJob(newJob, jobChannel) {
		logger.Error("channel is full: ", config.Setting.PROMETHEUS_CLIENT.QueueJobELements)
		return false
	}

	return true
}

// queueJob puts job into channel. If channel buffer is full, return false.
func queueJob(job JobObject, jobChan chan<- JobObject) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}

// worker processes jobs.
func worker(jobChan <-chan JobObject, wg *sync.WaitGroup) {
	// As soon as the current goroutine finishes (job done!), notify back WaitGroup.
	defer wg.Done()

	logger.Debug("Worker is waiting for jobs")

	for job := range jobChan {

		logger.Debug("Worker picked job: ", job)

		doPushJob(job)

		logger.Debug("Worker completed job", job)
	}
}

// worker processes jobs.
func doPushJob(jobElement JobObject) bool {

	logger.Debug("Do Sync Job", jobElement)
	var ok bool

	if jobElement.Type == GaugeType {
		var gaugeVec *prometheus.GaugeVec
		if gaugeVec, ok = config.Setting.PromGaugeMap[jobElement.NameStats]; ok {
			if err := prometheusPusher.Collector(gaugeVec).Push(); err != nil {
				logger.Error("Could not push gauge :["+jobElement.NameStats+"] time to Pushgateway:", err.Error())
			}
		}
	} else if jobElement.Type == HistogramType {
		var histogramVec *prometheus.HistogramVec
		if histogramVec, ok = config.Setting.PromHistogramMap[jobElement.NameStats]; !ok {
			if err := prometheusPusher.Collector(histogramVec).Push(); err != nil {
				logger.Error("Could not push histogram:["+jobElement.NameStats+"] time to Pushgateway:", err.Error())
			}
		}

	} else if jobElement.Type == CounterType {
		var counterVec *prometheus.CounterVec
		if counterVec, ok = config.Setting.PromCounterMap[jobElement.NameStats]; ok {
			if err := prometheusPusher.Collector(counterVec).Push(); err != nil {
				logger.Error("Could not push counter:["+jobElement.NameStats+"] time to Pushgateway:", err.Error())
			}
		}
	}

	return true
}
