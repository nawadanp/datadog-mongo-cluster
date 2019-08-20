package main

import (
	"flag"
	"log"
	"time"
)

type paramsMongo struct {
	MongoURI       string
	MongoTimeoutMS int
	Interval       time.Duration
	Cluster        string
}

type params struct {
	Mongo      paramsMongo
	DatadogKey string
}

const (
	appName       = "mongo-cluster-stats_datadog"
	dataDogAPIURL = "https://api.datadoghq.com/api/v1"
)

func parseFlags() params {
	flagMongoURI := flag.String("mongo-uri", "mongodb://127.0.0.1:27017", "MongoDB URI - Should be a mongos. Default: 'mongodb://127.0.0.1:27017'")
	flagMongoCluster := flag.String("cluster", "", "Cluster name, used to define the \"cluster\" label.")
	flagMongoTimeoutMS := flag.Int("mongo-timeout", 5000, "MongoDB Connection timeout in ms. Default: 5000ms")
	flagMongoInterval := flag.Int("interval", 600, "Metrics update interval in s. Default: 600s")
	flagDatadogKey := flag.String("dd-key", "", "Datadog Key.")

	flag.Parse()
	var config params
	config.Mongo.MongoURI = *flagMongoURI
	if *flagMongoCluster == "" {
		log.Fatalln("-cluster should be defined")
	}
	config.Mongo.Cluster = *flagMongoCluster
	config.Mongo.MongoTimeoutMS = *flagMongoTimeoutMS
	config.Mongo.Interval = time.Duration(*flagMongoInterval) * time.Second
	config.DatadogKey = *flagDatadogKey
	if *flagDatadogKey == "" {
		log.Fatalln("-dd-key should be defined")
	}
	return config
}

func main() {
	log.Println("Starting job")
	config := parseFlags()
	mongoSession := openMongoSession(config.Mongo)
	defer mongoSession.Close()

	colDatabases := mongoSession.DB("config").C("databases")

	var series series
	series.Series = append(series.Series, getDatabasesCount(colDatabases, config.Mongo.Cluster))
	pushToDatadog(config.DatadogKey, series)
}
