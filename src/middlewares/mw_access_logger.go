package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/shivanshkc/ledgerguard/src/logger"
	"github.com/shivanshkc/ledgerguard/src/utils/httputils"
)

// AccessLogger logs the details of the incoming requests and outgoing responses.
func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Capturing the entry time of the request to later log the total execution time.
		entryTime := time.Now()
		// Wrapping the writer within a custom writer for persisting statusCode.
		customWriter := &responseWriterWithCode{ResponseWriter: writer}

		// Logging the entry.
		logRequestEntry(request)
		// Releasing the control to following middlewares/handlers.
		next.ServeHTTP(customWriter, request)
		// Logging the exit.
		logRequestExit(customWriter, request, entryTime)
	})
}

// logRequestEntry logs the entry of the provided request.
func logRequestEntry(req *http.Request) {
	ctx, log := req.Context(), logger.Get()
	ctxData := httputils.GetReqCtx(ctx)

	// Getting client's and own IP address.
	clientIP, _ := httputils.GetClientIP(req)
	ownIP, _ := httputils.GetOwnIP()

	// Logging request entry.
	log.Info(ctx, &logger.Entry{
		Payload: "started request execution",
		Request: &logger.NetworkRequest{
			Protocol:    "http",
			ID:          ctxData.ID,
			Method:      req.Method,
			URL:         req.URL.Path,
			RequestSize: req.ContentLength,
			ServerIP:    ownIP,
			ClientIP:    clientIP,
		},
	})
}

// logRequestExit logs the exit (response details) of the provided request.
func logRequestExit(writer *responseWriterWithCode, req *http.Request, entry time.Time) {
	ctx, log := req.Context(), logger.Get()
	ctxData := httputils.GetReqCtx(ctx)

	// Getting client's and own IP address.
	clientIP, _ := httputils.GetClientIP(req)
	ownIP, _ := httputils.GetOwnIP()

	// Calculating response content length.
	var resContentLength int64
	if contentLength := writer.Header().Get("content-length"); contentLength != "" {
		resContentLength, _ = strconv.ParseInt(contentLength, 10, 64) // nolint:gomnd
	}

	// Logging response exit.
	log.Info(ctx, &logger.Entry{
		Timestamp: time.Now(),
		Payload:   "completed request execution",
		Request: &logger.NetworkRequest{
			Status:       writer.statusCode,
			Protocol:     "http",
			ID:           ctxData.ID,
			Method:       req.Method,
			URL:          req.URL.Path,
			RequestSize:  req.ContentLength,
			ResponseSize: resContentLength,
			Latency:      time.Since(entry),
			ServerIP:     ownIP,
			ClientIP:     clientIP,
		},
	})
}