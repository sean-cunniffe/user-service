echo "Cleaning test cache"
go clean -testcache
echo "Running go tests"
cd ./userservice
go test ./...
cd -