build-collector:
	docker build -f docker/collector/Dockerfile -t collector .

build-writer:
	docker build -f docker/writer/Dockerfile -t writer .

docker-influxdb:
	docker run -d --net=host --name=influxdb influxdb:1.4.2

connect-influxdb:
	docker exec -it influxdb influx
