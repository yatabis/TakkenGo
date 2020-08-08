package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"TakkenGo/line"
	"TakkenGo/scheduler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	s := scheduler.Init()
	defer s.Close()

	e.POST("/callback", line.Callback)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
