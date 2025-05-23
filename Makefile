# ---------------------------------------------------------------------------- #
#   ______ ______   ___   _____   _    _  _____  ______  ___  ______ ______    #
#   | ___ \| ___ \ / _ \ /  __ \ | |  | ||_   _||___  / / _ \ | ___ \|  _  \   #
#   | |_/ /| |_/ // /_\ \| /  \/ | |  | |  | |     / / / /_\ \| |_/ /| | | |   #
#   |    / | ___ \|  _  || |     | |/\| |  | |    / /  |  _  ||    / | | | |   #
#   | |\ \ | |_/ /| | | || \__/\ \  /\  / _| |_ ./ /___| | | || |\ \ | |/ /    #
#   \_| \_|\____/ \_| |_/ \____/  \/  \/  \___/ \_____/\_| |_/\_| \_||___/     #
#                                                                              #
# ---------------------------------------------------------------------------- #

# Color Definitions
NO_COLOR      = \033[0m
OK_COLOR      = \033[32;01m
ERROR_COLOR   = \033[31;01m
WARN_COLOR    = \033[33;01m

# Directories
APP_NAME      = rancher-rbac-wizard
BIN_DIR       = ./bin
GO_BUILD      = $(BIN_DIR)/$(APP_NAME)

# Main Targets
.PHONY: run serve run-ui build-ui build-ui-and-embed build-backend fmt create-cluster delete-cluster deploy-ingress-nginx clean

## Run the application in serve mode
run:
	@echo "$(OK_COLOR)==> Running the application...$(NO_COLOR)"
	go run . serve

## Serve the application (build UI and backend first)
serve: build-backend
	@echo "$(OK_COLOR)==> Starting the application...$(NO_COLOR)"
	go run . serve

## Run UI in development mode
run-ui:
	@echo "$(OK_COLOR)==> Running UI in development mode...$(NO_COLOR)"
	cd ui && npm run dev

## Build the UI
build-ui:
	@echo "$(OK_COLOR)==> Building UI...$(NO_COLOR)"
	cd ui && npm run build

copy-ui-artifacts:
	@echo "$(OK_COLOR)==> Embedding UI files into Go application...$(NO_COLOR)"
	cp -r ui/dist/ internal/embed/

## Build the Go backend and place the binary in bin directory
build-backend: copy-ui-artifacts
	@echo "$(OK_COLOR)==> Building Go backend...$(NO_COLOR)"
	mkdir -p $(BIN_DIR)
	go build -o $(GO_BUILD)

## Format Go code and tidy modules
fmt:
	@echo "$(OK_COLOR)==> Formatting Go code and tidying modules...$(NO_COLOR)"
	go fmt ./...
	go mod tidy

# Kubernetes Cluster Management
## Create a Kubernetes cluster using Kind
create-k8s-cluster:
	@echo "$(OK_COLOR)==> Creating Kubernetes cluster...$(NO_COLOR)"
	kind create cluster --config dev/kind.yaml

## Delete the Kubernetes cluster
delete-k8s-cluster:
	@echo "$(OK_COLOR)==> Deleting Kubernetes cluster...$(NO_COLOR)"
	kind delete cluster -n 'rbac-wizard-dev'

## Deploy NGINX Ingress controller to the cluster
deploy-ingress-nginx:
	@echo "$(OK_COLOR)==> Deploying ingress-nginx...$(NO_COLOR)"
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	@echo "$(OK_COLOR)==> Waiting for ingress-nginx to be ready...$(NO_COLOR)"
	kubectl wait --namespace ingress-nginx \
		--for=condition=ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=180s

# Cleanup
## Remove built binaries and cleanup
clean:
	@echo "$(OK_COLOR)==> Cleaning up build artifacts...$(NO_COLOR)"
	rm -rf $(BIN_DIR)
