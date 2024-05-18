serve:
	go run . serve

build-ui:
	cd ui && npm run build

build-ui-completely: build-ui
	statik -src=./ui/dist/ -dest=./internal/ -f
	rm -rf ./ui/dist

fmt:
	go fmt ./...
	go mod tidy
