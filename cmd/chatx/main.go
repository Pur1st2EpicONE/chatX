// Package main is the entry point of the application.
// It is responsible for bootstrapping and running the app.
package main

import "chatX/internal/app"

// main initializes the application and starts its execution lifecycle.
func main() {

	app.Boot().Run()

}
