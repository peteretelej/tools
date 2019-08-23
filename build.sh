#!/bin/bash
set -e

echo "Vetting, testing & building all Go binaries"

names=(paperlog get_secret slack)

for n in "${names[@]}"
do
    cd $n
    go mod download
    go vet ./... 
    go test ./...
    go build ./...
    cd ../
done 

echo "All binaries successfully build see directories: ${names[@]} "
