package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"TakkenGo/line"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/callback", line.Callback)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
