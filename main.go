package main

import (
	"fmt"
	"net/http"

	"github.com/JaineelVora08/librproto/controller"
	"github.com/JaineelVora08/librproto/router"
)

func main() {
	fmt.Println("STARTING")
	defer controller.Dbpool.Close()
	fmt.Println("Router building...")
	r := router.Router()
	fmt.Println("Listening on port 8000...")
	err := http.ListenAndServe(":8000", r)
	fmt.Println("Server has exited with:", err)
}
