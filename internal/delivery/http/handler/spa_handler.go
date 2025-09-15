package handler

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/akfaiz/go-vue-starter-kit/internal/config"
	"github.com/labstack/echo/v4"
)

type SPAHandler struct {
	cfg        config.Config
	root       fs.FS
	indexBytes []byte
	fileServer http.Handler
}

func NewSPAHandler(root fs.FS, cfg config.Config) *SPAHandler {
	indexBytes, err := fs.ReadFile(root, "index.html")
	if err != nil {
		indexBytes = []byte("index.html not found")
	}

	fileServer := http.FileServer(http.FS(root))

	return &SPAHandler{
		cfg:        cfg,
		root:       root,
		indexBytes: indexBytes,
		fileServer: fileServer,
	}
}

type AppConfig struct {
	ApiBaseUrl string `json:"apiBaseUrl"`
	Env        string `json:"env"`
}

func (h *SPAHandler) Env(c echo.Context) error {
	appConfig := AppConfig{
		ApiBaseUrl: h.cfg.App.ApiBaseURL,
		Env:        h.cfg.App.Env,
	}
	b, _ := json.Marshal(appConfig)
	w := c.Response().Writer
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	// Prevent caching so changing env vars takes effect instantly
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	_, _ = w.Write([]byte("window.__APP_CONFIG__ = "))
	_, _ = w.Write(b)
	_, _ = w.Write([]byte(";"))

	return nil
}

func (h *SPAHandler) View(c echo.Context) error {
	upath := c.Request().URL.Path

	// Root path → serve index.html
	if upath == "/" {
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := c.Response().Write(h.indexBytes)
		return err
	}

	// Check if file exists in FS
	fpath := strings.TrimPrefix(filepath.Clean(upath), "/")
	f, err := h.root.Open(fpath)
	if err == nil {
		_ = f.Close()
		h.setCacheHeaders(c.Response().Writer, fpath)
		h.fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	}

	// Fallback → index.html
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = c.Response().Write(h.indexBytes)
	return err
}

func (h *SPAHandler) setCacheHeaders(w http.ResponseWriter, p string) {
	ext := strings.ToLower(filepath.Ext(p))

	switch ext {
	case ".js", ".css", ".png", ".jpg", ".jpeg", ".webp", ".svg", ".ico", ".gif", ".woff2":
		if h.looksHashed(p) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
	default:
		w.Header().Set("Cache-Control", "no-cache")
	}
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
}

func (h *SPAHandler) looksHashed(p string) bool {
	base := filepath.Base(p)
	parts := strings.Split(base, ".")
	for _, s := range parts {
		if len(s) >= 8 {
			isHex := true
			for _, c := range s {
				if !strings.ContainsRune("0123456789abcdef", c) {
					isHex = false
					break
				}
			}
			if isHex {
				return true
			}
		}
	}
	return false
}
