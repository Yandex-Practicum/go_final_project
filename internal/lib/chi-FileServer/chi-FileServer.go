package chiFileServer

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"todo-list/internal/lib/common"

	"github.com/go-chi/chi"
)

const webPath = "./web"

func FileServer(r chi.Router, path string, root http.FileSystem) error {
	if strings.ContainsAny(path, "{}*") {
		return errors.New("fileServer does not permit any URL parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	return nil
}

func FileServerPath() (string, error) {

	appPath, err := common.ProjectRootPath()
	if err != nil {
		return "", fmt.Errorf("failed to get file server path: %w", err)
	}

	return filepath.Join(appPath, webPath), nil
}
