binary:
	cd cmd/hekate && \
	go build
	cd cmd/hctl && \
	go build
	cd cmd/debug-api && \
	go build
docker:
	docker build -t hekate:all-in-one -f build/allinone/Dockerfile .
	docker build -t hekate -f build/server/Dockerfile .
	docker build -t hekate-ui -f build/portal/Dockerfile .
run-windows:
	cd deployments && \
	run-windows.bat
apidocs:
	cd docs/api && \
	redoc-cli bundle api.yaml -o api.html
	cd docs/api && \
	redoc-cli bundle userapi.yaml -o userapi.html
