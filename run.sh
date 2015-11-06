#!/bin/bash

# Exit on error
set -e

# Test, build, then run with our sample data
go test
godebug build
./roster-generator.debug sample_players.txt sample_baggages.txt
