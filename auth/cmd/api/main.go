package main

import "fmt"

const version = "1.1.0" // set into .env

func main() {
	// Consume .env here
	cfg := config{
		addr: ":8080",
		env:  "development",
	}

	fmt.Println("Program start")
	app := application{
		config: cfg,
	}
	mux := app.mount()

	app.run(mux)
}
