#!/bin/bash -e

# go test
godep go test -v $(go list ./...  | grep -v "/try" | grep -v "/Godeps" | grep -v "/test/integration")

echo "Success!"
