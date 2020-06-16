package uploader

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var (
	osFs    = afero.NewOsFs() // afero is not needed
	isValid = func(info os.FileInfo) bool {
		return info.Size() < 1000000
	}
	// METHOD is the http method used by the request. default is POST
	METHOD = http.MethodPost
)

// UploadFolder ...
func UploadFolder(URL, fieldname, dirname string) (*http.Response, error) {

	// Creates a buffer to the body
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	// Walk each file in folder and create form part
	afero.Walk(osFs, dirname, walk(fieldname, w))
	w.Close()
	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := http.Client{}
	res, err := client.Do(req)
	return res, err
}

// walk will create a FormFile for each file in the folder
func walk(fieldname string, w *multipart.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		log.Println(path)
		if info.IsDir() {
			return nil
		}

		if !isValid(info) {
			return nil
		}
		file, err := osFs.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := w.CreateFormFile(fieldname, filepath.ToSlash(path))
		if err != nil {
			return err
		}

		io.Copy(part, file)

		return nil
	}

}

// SetIsValidFunction will override the default for the given function.
// This is function is used to filter the files that will be uploaded
func SetIsValidFunction(fn func(info os.FileInfo) bool) {
	isValid = fn
}
