build-collector:
	docker build -f docker/collector/Dockerfile -t collector .

build-writer:
	docker build -f docker/writer/Dockerfile -t writer .

build-evaluator:
	docker build -f - -t evaluator python/evaluator/ < docker/evaluator/Dockerfile

build-viewer:
	docker build -f - -t viewer python/viewer/ < docker/viewer/Dockerfile

docker-influxdb:
	docker run -d --net=host --name=influxdb influxdb:1.4.2

connect-influxdb:
	docker exec -it influxdb influx
