package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
}

type dockerPingResponse struct {
	OK     bool   `json:"ok"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type goVersionResponse struct {
	Version string `json:"version"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
	})

	mux.HandleFunc("/docker/ping", func(w http.ResponseWriter, r *http.Request) {
		output, err := runDockerVersion(r.Context())
		if err != nil {
			writeJSON(w, http.StatusOK, dockerPingResponse{
				OK:    false,
				Error: err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, dockerPingResponse{
			OK:     true,
			Output: output,
		})
	})

	mux.HandleFunc("/go/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, goVersionResponse{Version: runtime.Version()})
	})

	addr := ":" + getEnv("APP_PORT", "8080")
	log.Printf("go-sample app started on %s", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func runDockerVersion(ctx context.Context) (string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "docker", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(out))
		if trimmed != "" {
			return "", errWithOutput(err, trimmed)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

type outputError struct {
	err    error
	output string
}

func (e outputError) Error() string {
	return e.err.Error() + ": " + e.output
}

func errWithOutput(err error, output string) error {
	return outputError{err: err, output: output}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
