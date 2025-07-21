package main

import (
	_ "vk-internship/docs"
	"vk-internship/internal/app"
)

// @title VK Internship API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
