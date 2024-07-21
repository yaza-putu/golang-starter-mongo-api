package main

import (
	"runtime"

	"github.com/yaza-putu/golang-starter-mongo-api/internal/config"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/core"
)

func main() {
	// set max cpu
	runtime.GOMAXPROCS(config.App().MaxCpu)

	// load env
	core.Env()

	// init database
	core.Mongo()

	// redis
	core.Redis()
	// start server
	core.HttpServe()
}
