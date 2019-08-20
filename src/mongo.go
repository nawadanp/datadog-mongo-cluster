package main

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
)

type confDatabase struct {
	ID          string `bson:"_id"`
	Primary     string `bson:"primary"`
	Partitioned bool   `bson:"partitioned"`
}

func openMongoSession(config paramsMongo) *mgo.Session {
	session, err := mgo.DialWithInfo(buildDial(config))
	if err != nil {
		log.Fatalf("could not connect to %s: %s\n", config.MongoURI, err)
	}
	log.Printf("Connected to %s\n", config.MongoURI)

	return session
}

func buildDial(config paramsMongo) *mgo.DialInfo {
	dialInfo, err := mgo.ParseURL(config.MongoURI)
	if err != nil {
		log.Fatalf("Can't parse mongo URI : %s\n", err)
	}

	if dialInfo.AppName == "" {
		dialInfo.AppName = appName
	}

	dialInfo.Timeout = time.Duration(config.MongoTimeoutMS) * time.Millisecond
	dialInfo.ReadTimeout = 60 * time.Minute
	return dialInfo
}

func getDatabasesCount(colDatabases *mgo.Collection, cluster string) metric {
	// Define default metric settings
	var dbCount metric
	dbCount.Type = "gauge"
	dbCount.Metric = "custom.mongodb.cluster.databases"

	// Add metric tags
	tagCluster := "cluster:" + cluster
	dbCount.Tags = append(dbCount.Tags, tagCluster)

	// Get databases count
	dbTotal, err := colDatabases.Count()
	if err != nil {
		log.Fatalf("Error while counting number of database : %s\n", err)
	}

	if dbTotal == 0 {
		log.Fatalln("'databases' collection is empty or missing. Are you connected on a mongos ?")
	}

	var point points
	point[0] = float64(time.Now().Unix())
	point[1] = float64(dbTotal)

	dbCount.Points = append(dbCount.Points, point)
	return dbCount
}
