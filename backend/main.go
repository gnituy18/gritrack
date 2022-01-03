package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/valyala/fasthttp"

	"gritter/pkg/api"
	"gritter/pkg/log"
)

func main() {
	if err := fasthttp.ListenAndServe(":"+os.Getenv("PORT"), api.Router().HandleRequest); err != nil {
		log.Global().Error("fasthttp.ListenAndServe failed in main")
	}
}
