#!/bin/bash

echo "Running deploy ..."
pwd

# build bood
go build ./build/cmd/bood

# build & test bood with new bood
cd ./build
../bood

# build example
cd ../example
../bood

