package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/xid"
	"github.com/spf13/afero"
)

var osFs = afero.NewOsFs()

const root = "examples/uploadfolder/server"

func upload(c echo.Context) error {
	reqID := xid.New().String()

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}

		// generate the filepath and and creates all the necessary nested directories
		baseDIR := fmt.Sprintf("%s/requests/%s/", root, reqID)              // base directory for the request "requests/reqID/"
		reqDIR := fmt.Sprintf("%s%s", baseDIR, filepath.Dir(file.Filename)) // "requests/reqID/dir1/dir2/"
		log.Println(reqDIR)
		err = osFs.MkdirAll(reqDIR, os.ModePerm)
		if err != nil {
			return err
		}

		// Destination
		dst, err := os.Create(baseDIR + file.Filename) // "requests/reqID/dir1/dir2/file.txt"
		if err != nil {
			return err
		}

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		src.Close()
		dst.Close()
	}

	return c.String(http.StatusOK, fmt.Sprintf("uploaded successfully %d files", len(files)))
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", root+"/public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":5000"))
}
