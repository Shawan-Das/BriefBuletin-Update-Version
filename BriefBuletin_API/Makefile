.PHONY: all build generate tools quickbuild

BINARY=briefbuletin
TOOL=hrmstools
OUTPUT=build/exec
MAIN=main.go  # Corrected path
DCONFIG := config/connection/dev-config.json 

wbr: winBuild runOnWindows 

server:
	go mod download
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-w -extldflags '-static' " -o ${OUTPUT}/${BINARY} 
	
sqlc:
	sqlc generate --file config/sqlc/db_query/db.yaml

swag:
	swag init -g main.go --output docs

# Windows build
winBuild:
	go mod download
	go mod tidy
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}

# Use 'make runOnWindows' to execute program
# $(MAKE) winBuild
dev:
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}
	./${OUTPUT}/${BINARY}.exe -c ./config/connection/dev-config.json --port 7070

run:
	./${OUTPUT}/${BINARY} -c ./config/connection/dev-config.json --port 7070
