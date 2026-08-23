package web

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed all:dist
var distFS embed.FS

func DistFS() http.FileSystem {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}

func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}

func RegisterRoutes(e *echo.Echo) {
	embeddedFS := FS()

	serveEmbeddedFile := func(c echo.Context, filePath string) error {
		f, err := embeddedFS.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return fs.ErrInvalid
		}

		var rs io.ReadSeeker
		if seeker, ok := f.(io.ReadSeeker); ok {
			rs = seeker
		} else {
			data, readErr := io.ReadAll(f)
			if readErr != nil {
				return readErr
			}
			rs = bytes.NewReader(data)
		}

		http.ServeContent(c.Response(), c.Request(), stat.Name(), stat.ModTime(), rs)
		return nil
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet && c.Request().Method != http.MethodHead {
				return next(c)
			}

			reqPath := c.Request().URL.Path
			if strings.HasPrefix(reqPath, "/api") || reqPath == "/openapi.json" || reqPath == "/health" {
				return next(c)
			}

			clean := path.Clean(reqPath)
			trimmed := strings.TrimPrefix(clean, "/")

			if trimmed == "" || trimmed == "." {
				if err := serveEmbeddedFile(c, "index.html"); err == nil {
					return nil
				}
			}

			if err := serveEmbeddedFile(c, trimmed); err == nil {
				return nil
			}

			if err := serveEmbeddedFile(c, trimmed+".html"); err == nil {
				return nil
			}

			if err := serveEmbeddedFile(c, path.Join(trimmed, "index.html")); err == nil {
				return nil
			}

			if err := serveEmbeddedFile(c, "index.html"); err == nil {
				return nil
			}

			return next(c)
		}
	})
}
