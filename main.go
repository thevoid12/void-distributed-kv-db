package main

import (
	"flag"
	"log"
	"main/config"
	"main/handler"
	"main/pkg/db"
	"net/http"

	"github.com/spf13/viper"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	shard      = flag.String("shard", "", "The name of the shard for the data")
	configFile = flag.String("config-file", "config/config.json", "Config file for static sharding")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatalf("Must provide db-location")
	}
	if *shard == "" {
		log.Fatalf("Must provide shard name")
	}

}

func main() {
	parseFlags()

	log.Println("Reading configuration from:", *configFile)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("config/") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var shardConfig []config.Shard
	err = viper.UnmarshalKey("shards", &shardConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling shard config: %v", err)
	}

	shards, err := config.ParseShards(shardConfig, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config: %v", err)
	}

	log.Printf("Shard count is %d, current shard: %d", shards.Count, shards.CurIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
	defer close()

	srv := handler.NewServer(db, shards)

	// Register handlers
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Printf("Starting server on %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
