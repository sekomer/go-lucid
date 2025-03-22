#!/bin/bash

DEV_CLUSTER_SIZE=1

pids=()

cleanup() {
    echo "Shutting down all nodes..."
    for pid in "${pids[@]}"; do
        kill $pid 2>/dev/null
    done
    exit 0
}

# set up trap to catch SIGINT (Ctrl+C) and SIGTERM
trap cleanup SIGINT SIGTERM

go run main.go --config config/bootnode.yaml node start --boot &
pids+=($!)

sleep 1

for i in {1..$DEV_CLUSTER_SIZE}; do
    go run main.go --config config/devnode.yaml node start --dev &
    pids+=($!)
done

wait