docker-build:
	docker build -f docker/collector/Dockerfile -t collector .

docker-influxdb:
	docker run -d --net=host --name=influxdb influxdb:1.4.2

connect-influxdb:
	docker exec -it influxdb influx
