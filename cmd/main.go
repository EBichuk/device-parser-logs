package main

import (
	"device-parser-logs/internal/api"
)

func main() {
	api := api.New()
	api.RunApi()
}
