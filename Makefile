build-docker:
	docker build -t jwt-server -f build/server/Dockerfile .
	docker build -t jwt-server-portal -f build/portal/Dockerfile .
