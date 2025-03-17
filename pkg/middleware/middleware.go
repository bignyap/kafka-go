package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(mux http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}
	return mux
}
