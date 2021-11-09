package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/valyala/fasthttp"

	"gritter/pkg/api"
	"gritter/pkg/log"
)

func main() {
	if err := fasthttp.ListenAndServe(":8080", api.Router().HandleRequest); err != nil {
		log.Global().Error("fasthttp.ListenAndServe failed in main")
	}
}
