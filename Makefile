.PHONY: all build run 

# Default target that runs both build and run
all: build run

build:
# Change directory to the location of this script
	cd $(dirname $0)

#building the executable and using -o to name it as main
	go build -o main

run:
# set -e is used to exit immediately if any command fails
	set -e

# trap 'killall main' SIGINT sets up a "trap" that listens for the SIGINT signal  triggered by Ctrl+c.
# When this signal is detected, it will execute killall distribkv,
# which terminates all running instances of the distribkv process
	trap 'killall main' SIGINT

# Change directory to the location of this script
	cd $(dirname $0)

# Kill any existing distribkv processes. I added a || true to make sure that
# since killall might return a non-zero exit code if there are no processes to kill  false || true will always return true 
# thus we are not exiting the program
	killall main || true

# Pause for 0.1 seconds to ensure processes are terminated
	sleep 0.1

# Start mukliple instances(3) of void-distributed-kv-db with different configurations
	./main -db-location=voidZoneA.db -http-addr=127.0.0.1:8080 -config-file=config/config.json -shard=void-shard-zoneA &
	./main -db-location=voidZoneB.db -http-addr=127.0.0.1:8081 -config-file=config/config.json -shard=void-shard-zoneB &
	./main -db-location=voidZoneC.db -http-addr=127.0.0.1:8082 -config-file=config/config.json -shard=void-shard-zoneC &
	wait
