build:
	python3 -u build-config.py
	cd node-exporter-proxy && docker build -t songjiayang/node-exporter-proxy:0.1.0 .
start:
	docker-compose up -d 
stop:
	docker-compose down