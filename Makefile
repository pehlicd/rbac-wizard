serve: build-ui-completely
	go run . serve

run-ui:
	cd ui && npm run dev

build-ui:
	cd ui && npm run build

build-ui-completely: build-ui
	statik -src=./ui/dist/ -dest=./internal/ -f
	rm -rf ./ui/dist

fmt:
	go fmt ./...
	go mod tidy
