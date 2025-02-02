SHA = $(shell git rev-parse HEAD)

build:
	GOOS=linux GOARCH=amd64 go build zpool_status.go

clean:
	rm -rf ./zpool_status

docker:
	docker build -t skirsch10/zpool_status:$(SHA) --platform=linux/amd64 .

push:
	docker push skirsch10/zpool_status:$(SHA)