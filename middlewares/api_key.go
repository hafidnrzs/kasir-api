package middlewares

import "net/http"

// func (api key) func handler http.handler
// fungsi middleware ini me-return fungsi handler. Fungsi handler tersebut juga me-return fungsi
func APIKeyMiddleware(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-Api-Key")

			if apiKey == "" {
				http.Error(w, "API Key required", http.StatusUnauthorized)
				return
			}

			if apiKey != validApiKey {
				http.Error(w, "Invalid API Key", http.StatusUnauthorized)
				return
			}

			// jalankan fungsi selanjutnya jika API Key valid
			next(w, r)
		}
	}
}
