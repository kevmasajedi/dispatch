package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/db"
	"strings"
	"os"
)

func fetch_domain_worker_url(query map[string]string) string {
	if db.Connect() {
		res := db.FindOneFromCollection("domain_workers", query)
		if res == nil {
			return ""
		}

		ip, ipOk := res["ip"]
		port, portOk := res["port"]
		if !ipOk || !portOk {
			return ""
		}
		return fmt.Sprintf("http://%s:%s", ip, port)
	}
	return ""
}
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 || parts[1] == "" {
		http.Error(w, "Dispatch Works!", http.StatusBadRequest)
		return
	}
	domainName := parts[1]

	q := make(map[string]string)
	q["domain"] = domainName

	domain_worker_url := fetch_domain_worker_url(q)
	if domain_worker_url == "" {
		http.Error(w, "Incomplete domain worker data", http.StatusInternalServerError)
	}

	proxyUrl, err := url.Parse(domain_worker_url)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	r.URL.Path = strings.Join(parts[2:], "/")
	proxy.ServeHTTP(w, r)
}
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        	enableCors(&w) // Add CORS headers
        	ProxyHandler(w, r) // Handle the request
    	})
	port := "8012"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Printf("Dispatch is running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
