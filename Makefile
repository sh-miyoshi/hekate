build-docker:
	docker build -t hekate -f build/server/Dockerfile .
	docker build -t hekate-ui -f build/portal/Dockerfile .
