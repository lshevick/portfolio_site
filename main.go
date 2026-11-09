package main
import (
    "log"
    "net/http"
    "os"
)

func main() {
    siteDir := "/home/lshevick/portfolio/main"
    
    if _, err := os.Stat(siteDir); os.IsNotExist(err) {
        log.Fatalf("frontend directory not found: %s", siteDir)
    }

    mux := http.NewServeMux()

    // Serve static files from /main
    fileServer := http.FileServer(http.Dir(siteDir))
    
    // Health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    // Serve the root path with fileServer
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/health" {
            return
        }
        fileServer.ServeHTTP(w, r)
    })

    addr := os.Getenv("PORT")
    if addr == "" {
        addr = "8080"
    }

    log.Printf("serving frontend from %s on :%s", siteDir, addr)
    log.Fatal(http.ListenAndServe(":"+addr, mux))
}