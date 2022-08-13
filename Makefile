SHELL := /bin/bash

run_server:
	dapr run --app-id hello-dapr -p 8088 -H 8089 -G 3500 --config ./pipeline.yaml --components-path ./components/ go run main.go

.PHONY: compile