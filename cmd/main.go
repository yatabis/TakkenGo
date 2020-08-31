package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/yatabis/TakkenGo/line"
	"github.com/yatabis/TakkenGo/scheduler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	s := scheduler.Init()
	defer s.Close()

	e.GET("/ping", ping)
	e.POST("/callback", line.Callback)
	e.GET("/training", training)
	e.GET("/snooze", snooze)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func training(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "shortcuts://run-shortcut?name=takken-go")
}

func snooze(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "shortcuts://run-shortcut?name=takken-go/snooze")
}
