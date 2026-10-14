// Command duanju_core_server is the standalone HTTP server that exposes the
// native core (duanjuapp/native/core) to local clients — the Kodi plugin
// (plugin.video.hongguo, see resources/lib/zjcore.py) and the Android APK.
//
// Wire protocol (must match zjcore.py):
//
//	GET  /healthz                     -> 200 "ok"
//	POST /  (application/json action) -> core.NativeRequest(json) -> json
//
// It binds 127.0.0.1:${DUANJU_CORE_PORT:-18765}. The proxy is supplied by the
// caller through the standard HTTP_PROXY/HTTPS_PROXY environment variables;
// the core reads it from the environment, so this wrapper needs no proxy logic.
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"duanjuapp/native/core"
)

func main() {
	port := os.Getenv("DUANJU_CORE_PORT")
	if port == "" {
		port = "18765"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"ok":false,"error":"读取请求体失败"}`)
			return
		}
		out := core.NativeRequest(string(body))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, out)
	})

	addr := "127.0.0.1:" + port
	log.Printf("duanju_core_server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("duanju_core_server exited: %v", err)
	}
}
