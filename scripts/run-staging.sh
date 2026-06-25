#!/bin/sh
set -e
# Initialize Swagger documentation
# go install github.com/swaggo/swag/cmd/swag@latest
# swag init


# Run the application using air for live-reloading
air -c .air.toml
air ./cmd/main.go -b 0.0.0.0