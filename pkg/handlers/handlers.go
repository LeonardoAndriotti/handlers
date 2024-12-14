package handlers

import (
	"fmt"
	"net/http"
)

type RestClient interface {
	POST(pattern string, handler func(http.ResponseWriter, *http.Request))
	GET(pattern string, handler func(http.ResponseWriter, *http.Request))
	Build(addr ...string)
}

type handler struct {
	port string
}

func NewHandler(vars ...map[string]string) RestClient {

	return &handler{
		port: ":8080",
	}
}

func (c *handler) POST(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.Handle(pattern, methodValidatorMiddleware(http.MethodPost, http.HandlerFunc(handler)))
}

func (c *handler) GET(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.Handle(pattern, methodValidatorMiddleware(http.MethodGet, http.HandlerFunc(handler)))
}

func (c *handler) Build(addr ...string) {
	port := ":" + c.port
	if len(addr) > 0 {
		port = addr[0]
	}

	err := http.ListenAndServe(port, nil)
	if err != nil {
		return
	}
}

// Middleware para validar o método HTTP
func methodValidatorMiddleware(allowedMethod string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configura os cabeçalhos CORS
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Se for uma requisição OPTIONS, finalize aqui
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != allowedMethod {
			http.Error(w, fmt.Sprintf("Método não permitido. Use %s", allowedMethod), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
