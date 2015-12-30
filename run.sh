#!/bin/bash

# Exit on error
set -e

# Test, build, then run with our sample data
go test
go build
./roster-generator sample_players.csv sample_baggages.csv
