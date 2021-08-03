package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AccessLogMiddleware(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			fields := logrus.Fields{
				"type":        "access log",
				"method":      r.Method,
				"remote_addr": r.RemoteAddr,
				"host":        r.Host,
			}
			logger.WithFields(fields).Info()
		})
	}
}

func TimeLogMiddleware(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			fields := logrus.Fields{
				"type":      "time log",
				"work_time": duration,
			}
			logger.WithFields(fields).Info()
		})
	}
}

const CorrelationID = "correlation_id"

func IDMiddleware(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id := uuid.New()
			//revive:disable
			ctx = context.WithValue(ctx, CorrelationID, id.String()) //nolint:staticcheck
			//revive:enable
			r = r.WithContext(ctx)
			logger.WithFields(logrus.Fields{
				CorrelationID: id,
			}).Info()
			next.ServeHTTP(w, r)
		})
	}
}
