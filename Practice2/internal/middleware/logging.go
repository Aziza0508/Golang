package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006-01-02T15:04:05")
		method := strings.ToUpper(r.Method[:1]) + strings.ToLower(r.Method[1:])
		fmt.Fprintf(os.Stdout, "%s %s %s request received\n", timestamp, method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
