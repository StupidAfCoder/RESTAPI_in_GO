package middlewares

import (
	"net/http"
)

var allowedOrgins = []string{
	"https://localhost:3000",
	"https://www.frontend.com",
}

// Here Cors stand for corss Origin resource sharing which helps to understand what origins are allowed to access the api here we will use an exmaple that our Api can be accesed by two domains
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// No Origin header = same-origin request, just pass through
			next.ServeHTTP(w, r)
			return
		}
		if !isOriginAllowed(origin) {
			http.Error(w, "Not allowed by CORS policy", http.StatusForbidden)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrgins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
