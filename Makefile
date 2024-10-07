NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

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

create-cluster:
	@echo "$(OK_COLOR)==> Creating the k8s cluster$(NO_COLOR)"
	@kind create cluster --config dev/kind.yaml

delete-cluster:
	@echo "$(OK_COLOR)==> Deleting the k8s cluster$(NO_COLOR)"
	@kind delete cluster -n 'rbac-wizard-dev'

deploy-ingress-nginx:
	@echo "$(OK_COLOR)==> Deploying the ingress-nginx$(NO_COLOR)"
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	@kubectl wait --namespace ingress-nginx \
		--for=condition=ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=180s
		