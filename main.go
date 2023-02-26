package main

import (
	"fmt"
)

func main() {
	//first init() function executes that present in init.go
	fmt.Println("App initialization is done.")

	fmt.Println("server is starting...")
	app() //it init server and business logic(service layre)
}
