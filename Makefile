SHELL := /bin/bash

run_server:
	dapr run --app-id hello-dapr -p 8088 -H 8089 -G 3500 --config ./pineline.yaml --components-path ./components/ go run main.go

run_zipkin:
	dapr run --app-id hello-dapr -p 8088 -H 8089 -G 3500 --config ./tracing.yml --components-path ./components/ go run main.go

.PHONY: compile