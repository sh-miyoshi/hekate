build-docker:
	docker build -t hekate:all-in-one -f build/allinone/Dockerfile .
	docker build -t hekate -f build/server/Dockerfile .
	docker build -t hekate-ui -f build/portal/Dockerfile .
run-windows:
	cd cmd/hekate && \
	go build && \
	${env}:HEKATE_PORTAL_ADDR = "http://localhost:3000" && \
	start hekate.exe
	cd cmd/portal && \
	npm run dev