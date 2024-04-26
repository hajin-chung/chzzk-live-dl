package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s\n", err)
	}

	InitClient()

	streamers := Streamers{}
	err = streamers.Load()
	if err != nil {
		log.Fatalf("error loading streamer from file: %s\n", err)
	}
	streamers.Watch()

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/api/streamer", func(c echo.Context) error {
		return c.JSON(200, streamers.Infos)
	})

	e.GET("/api/streamer/search", func(c echo.Context) error {
		query := c.QueryParam("query")
		log.Println(query)
		data, err := SearchChannel(query)
		if err != nil {
			return err
		}

		return c.JSON(200, data)
	})

	e.POST("/api/streamer/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := streamers.AddStreamer(id)
		if err != nil {
			return err
		}

		return c.JSON(200, streamers.Infos)
	})

	e.DELETE("/api/streamer/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := streamers.DeleteStreamer(id)
		if err != nil {
			return err
		}

		return c.JSON(200, streamers.Infos)
	})

	e.POST("/api/streamer/:id/download", func(c echo.Context) error {
		id := c.Param("id")
		err := streamers.StartDownload(id)
		if err != nil {
			return err
		}

		return c.JSON(200, streamers.Infos)
	})

	e.DELETE("/api/streamer/:id/download", func(c echo.Context) error {
		id := c.Param("id")
		err := streamers.StopDownload(id)
		if err != nil {
			return err
		}

		return c.JSON(200, streamers.Infos)
	})

	e.POST("/api/streamer/:id/autoDownload", func(c echo.Context) error {
		id := c.Param("id")
		streamers.Infos[id].AutoDownload = true

		return c.JSON(200, streamers.Infos)
	})

	e.DELETE("/api/streamer/:id/autoDownload", func(c echo.Context) error {
		id := c.Param("id")
		streamers.Infos[id].AutoDownload = false

		return c.JSON(200, streamers.Infos)
	})

	e.POST("/api/cred", func(c echo.Context) error {
		bytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		credentials := &Credentials{}
		err = json.Unmarshal(bytes, credentials)
		if err != nil {
			return err
		}

		os.Setenv("NID_AUT", credentials.NID_AUT)
		os.Setenv("NID_SES", credentials.NID_SES)

		return c.JSON(200, credentials)
	})

	e.Logger.Fatal(e.Start(":2000"))
}
