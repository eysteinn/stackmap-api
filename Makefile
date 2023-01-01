BUILDVERSION:=latest

get-docs:
	go get -u github.com/swaggo/swag/cmd/swag

docs: get-docs
	swag init --dir cmd/api --parseDependency --output docs

build:
	go build -o bin/restapi cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./test/...

build-docker: build
	docker build . -t api-rest:$(BUILDVERSION)

run-docker: build-docker
	docker run -p 3000:3000 api-rest:$(BUILDVERSION)

port-forward:
	kubectl port-forward svc/postgresql 5432:5432

kind-deploy: build-docker
	kind load docker-image api-rest:$(BUILDVERSION) && kubectl apply -f deployment.yaml 

#swagger-build:
#	swagger generate spec -i ./swagger/swagger_base.yaml -o ./swagger.yaml

swagger-serve:
	cd swagger && swagger serve --flatten --port=9009 -F=swagger swagger.yaml
