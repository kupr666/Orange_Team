package core_http_middleware

import (
	"net/http"
	
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

  func RequestID() Middleware {
      return func(next http.Handler) http.Handler {
          return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
              requestID := r.Header.Get(requestIDHeader)
              if requestID == "" {
                  requestID = uuid.NewString()
              }

              w.Header().Set(requestIDHeader, requestID)
              next.ServeHTTP(w, r)
          })
      }
  }
