.DEFAULT_GOAL := swagger

install_swagger:
	go get -u -d github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....

	go run github.com/go-swagger/go-swagger/cmd/swagger generate spec -o ./swagger.yaml --scan-models
