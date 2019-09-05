package handler

import (
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/index.html", http.StatusFound)
		return
	}
}
