package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

func uploadTemplate(message string) string {
	return `<!DOCTYPE html>
               <html>
                   <head></head>
                   <body>
                       <form enctype="multipart/form-data" action="/upload" method="POST">
                           <h3>File Upload</h3>
                           <input type="file" placeholder="uploadfile" name="uploadfile"><br>
                           <input type="submit" value="Upload">
                           <div>` + message + `</div>
                       </form>
                   <body>
               </html>`
}

func formatApiError(err error) string {
	return fmt.Sprintf(`{"status":"error", "error": "%v"}`, err.Error())
}

func newAssetId() string {
	return fmt.Sprintf("%s", uuid.Must(uuid.NewV4()))
}

func FileUpload(r *http.Request) (string, int) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if nil != err {
		logger.Error(err)
		return formatApiError(err), http.StatusBadRequest
	}

	defer file.Close()

	asset_id := newAssetId()
	ext := filepath.Ext(handler.Filename)
	save_file_name := fmt.Sprintf("%v%v", asset_id, ext)

	f, err := os.OpenFile(ASSETS_DIRECTORY+save_file_name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		return formatApiError(err), http.StatusInternalServerError
	}
	defer f.Close()
	io.Copy(f, file)

	err = Insert(handler.Filename, asset_id)
	if nil != err {
		logger.Error(err)
		return formatApiError(err), http.StatusInternalServerError
	}

	result, err := Select(asset_id)
	if nil != err {
		logger.Error(err)
		return formatApiError(err), http.StatusInternalServerError
	}

	return fmt.Sprintf(`{"status":"ok", "data": %v}`, result), http.StatusOK
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		fmt.Fprintf(w, uploadTemplate(""))
		return
	}

	results, status_code := FileUpload(r)
	w.WriteHeader(status_code)
	fmt.Fprintf(w, uploadTemplate(results))
}

func FileUploadApiV1Handler(w http.ResponseWriter, r *http.Request) {
	results, status_code := FileUpload(r)
	w.WriteHeader(status_code)
	fmt.Fprintf(w, results)
}

// TODO
//  - DELETE
