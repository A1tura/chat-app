package errors

import (
	"fmt"
	"net/http"
)

func NotFoundError(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Error 404: Not Found")
}
