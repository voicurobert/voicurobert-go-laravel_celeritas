package middleware

import "net/http"

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !m.App.Session.Exists(request.Context(), "userID") {
			http.Error(writer, http.StatusText(401), http.StatusUnauthorized)
		}
	})
}
