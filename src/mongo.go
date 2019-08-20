package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
)

type confDatabase struct {
	ID          string `bson:"_id"`
	Primary     string `bson:"primary"`
	Partitioned bool   `bson:"partitioned"`
}

func openMongoSession(config paramsMongo) (*mgo.Session, error) {
	mongoDialInfo, err := buildDial(config)
	if err != nil {
		e := fmt.Errorf("unable to build dialinfo: %s", err)
		return nil, e
	}
	session, err := mgo.DialWithInfo(mongoDialInfo)
	if err != nil {
		e := fmt.Errorf("could not connect to %s: %s", config.MongoURI, err)
		return nil, e
	}
	log.Printf("Connected to %s\n", config.MongoURI)

	return session, nil
}

func buildDial(config paramsMongo) (*mgo.DialInfo, error) {
	dialInfo, err := mgo.ParseURL(config.MongoURI)
	if err != nil {
		e := fmt.Errorf("can't parse mongo URI: %s", err)
		return nil, e
	}

	if dialInfo.AppName == "" {
		dialInfo.AppName = appName
	}

	dialInfo.Timeout = time.Duration(config.MongoTimeoutMS) * time.Millisecond
	dialInfo.ReadTimeout = 60 * time.Minute
	return dialInfo, nil
}

func getDatabasesCount(colDatabases *mgo.Collection, cluster string) (*metric, error) {
	// Define default metric settings
	//var dbCount *metric
	dbCount := &metric{}
	dbCount.Type = "gauge"
	dbCount.Metric = "custom.mongodb.cluster.databases"

	// Add metric tags
	tagCluster := "cluster:" + cluster
	dbCount.Tags = append(dbCount.Tags, tagCluster)

	// Get databases count
	dbTotal, err := colDatabases.Count()
	if err != nil {
		e := fmt.Errorf("Error while counting number of database: %s", err)
		return nil, e
	}

	if dbTotal == 0 {
		return nil, errors.New("'databases' collection is empty or missing. Are you connected on a mongos ?")
	}

	var point points
	point[0] = float64(time.Now().Unix())
	point[1] = float64(dbTotal)

	dbCount.Points = append(dbCount.Points, point)
	return dbCount, nil
}
