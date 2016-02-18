Target: help ## Description
	@echo

app=go-config
image=rmohid/$(app)

.FORCE:

help: ##  This help dialog
	@cat $(MAKEFILE_LIST) | perl -ne 's/(^\S+): .*##\s*(.+)/printf "\n %-16s %s", $$1,$$2/eg'; echo

build: .FORCE  ## Build the docker image
	GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(app)
	docker build -t $(image) .

run: stop ## Run service
	docker run -d -p 7100:7100 --name $(app)  $(image)

stop: ## Stop service
	docker  rm -f $(app) || true

inter: ## Run interactive test session
	docker run -it $(image) sh

push: ## Push image to docker hub
	docker push $(image)

test: ## Show the help index
	curl $$(docker-machine ip default):7100
	
clean: ## Remove app container and image
	docker  rm -f $(app) || true
	docker images | awk '/$(app)/{print $$3}' | xargs docker rmi -f

clean_all: ## Remove all app images and exited containers
	docker rmi $$(docker images -f "dangling=true" -q) || true
	docker rm -v $$(docker ps -a -q -f status=exited) || true
