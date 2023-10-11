BUILDVERSION:=latest
DOCKERIMAGE:=stackmap-api:$(BUILDVERSION)

#get-docs:
#	go get -u github.com/swaggo/swag/cmd/swag

#docs: get-docs
#	swag init --dir cmd/api --parseDependency --output docs
.PHONY: kind-load
.PHONY: build-docker

build:
	go build -o bin/restapi cmd/api/main.go

run:
	go run cmd/api/main.go

#test:
#	go test -v ./test/...

build-docker: #build
	docker build . -t $(DOCKERIMAGE)

run-docker: build-docker
	docker run -p 3000:3000 $(DOCKERIMAGE)

port-forward:
	kubectl port-forward svc/postgresql 5432:5432

kind-load: build-docker
	kind load docker-image $(DOCKERIMAGE)

kind-deploy: kind-load
	kubectl apply -f deployment.yaml 

k3s-deploy:
	docker save $(DOCKERIMAGE) | sudo k3s ctr images import -

k-delete:
	kubectl delete deploy stackmap-api

k-deploy:
	kubectl apply -f https://raw.githubusercontent.com/eysteinn/stackmap-api/main/deployment.yaml

#swagger-build:
#	swagger generate spec -i ./swagger/swagger_base.yaml -o ./swagger.yaml

swagger-serve:
	cd swagger && swagger serve --flatten --port=9009 -F=swagger swagger.yaml
