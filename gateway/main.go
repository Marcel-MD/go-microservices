package main

import (
	"gateway/http"
)

func main() {

	srv := http.GetServer()
	srv.Run()

}
