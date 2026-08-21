package openapi

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/fincher/main.go -d ../ --outputTypes json -o . --parseInternal --parseDepth 5
