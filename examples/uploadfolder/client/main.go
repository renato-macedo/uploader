package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/renato-macedo/uploader"
)

const root = "examples/uploadfolder/client"

func main() {
	// endpoint that will handle multipart requests
	res, err := uploader.UploadFolder("http://localhost:5000/upload", "files", root+"/files")
	if err != nil {
		log.Printf("Upload error %s", err.Error())
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	fmt.Printf("Response %s", b)
}
