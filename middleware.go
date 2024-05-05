package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Context keys.
type contextKey int

const (
	requestIDKey contextKey = iota
	traceIDKey
)

// Logger setup
var logger = logrus.New()

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()   // Generate a new UUID as TraceID
		requestID := uuid.New().String() // Generate a new UUID as RequestID

		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		ctx = context.WithValue(ctx, requestIDKey, requestID)

		// Log the IDs using Logrus
		logger.WithFields(logrus.Fields{
			"traceID":   traceID,
			"requestID": requestID,
		}).Info("New request")

		// Adding the TraceID to the response header for debugging/tracing
		w.Header().Add("X-Trace-ID", traceID)
		w.Header().Add("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
