package main

import (
	"io"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func server() {
	e := echo.New()
	e.GET("file.wav", func(c echo.Context) error {
		cap, err := CaptureSpeaker("alsa_output.pci-0000_00_1b.0.analog-stereo.monitor")
		if err != nil {
			panic(err)
		}
		defer cap.Close()
		c.Response().Header().Set(echo.HeaderContentType, "audio/x-wav")
		c.Response().WriteHeader(http.StatusOK)

		io.Copy(c.Response(), cap.Buffer)

		return nil
	})
	e.Logger.Fatal(e.Start(":1323"))

}
