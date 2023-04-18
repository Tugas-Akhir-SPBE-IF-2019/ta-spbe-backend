package assessment

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type M map[string]interface{}

func (handler *assessmentHandler) DownloadSupportDocuments(w http.ResponseWriter, r *http.Request) {
	files := []map[string]interface{}{}
	basePath, _ := os.Getwd()
	filesLocation := filepath.Join(basePath, "static/supporting-documents")

	var pathFile string
	err := filepath.Walk(filesLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		pathFile = path
		files = append(files, map[string]interface{}{"filename": info.Name(), "path": path})
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// res, err := json.Marshal(files)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(res)
	f, err := os.Open(pathFile)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
