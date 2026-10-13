package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogHTTPMiddleware(hiddenHeaders []string) echo.MiddlewareFunc {
	var hidden = append([]string{}, hiddenHeaders...)
	hidden = append(hidden, "authorization")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ts      = time.Now()
				traceID = uuid.New().String()
				header  = c.Request().Header.Clone()
			)

			req := c.Request()
			ctx := context.WithValue(req.Context(), "traceID", traceID)
			c.SetRequest(req.WithContext(ctx))

			for _, item := range hidden {
				if v := header.Get(item); len(v) != 0 {
					header.Set(strings.ToLower(item), "*")
				}
			}

			var reqBodyBuf bytes.Buffer
			if c.Request().Body != nil {
				bodyBytes, _ := io.ReadAll(c.Request().Body)
				reqBodyBuf.Write(bodyBytes)
				c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			logrus.
				WithField("http.header", header).
				WithField("method", c.Request().Method).
				WithField("traceID", traceID).
				WithField("path", c.Request().URL.Path).
				WithField("query", c.Request().URL.RawQuery).
				WithField("request.body", reqBodyBuf.String()).
				Infof("incoming request from client")

			defer func() {
				if r := recover(); r != nil {
					logrus.
						WithField("traceID", traceID).
						WithField("duration", time.Since(ts).String()).
						Errorf("panic: %v", r)
					_ = c.NoContent(http.StatusInternalServerError)
				}
			}()

			origWriter := c.Response().Writer
			bodyWriter := &bodyCaptureWriter{ResponseWriter: origWriter, body: new(bytes.Buffer)}
			c.Response().Writer = bodyWriter

			err := next(c)
			duration := time.Since(ts)

			respBody := bodyWriter.body.String()
			if err != nil || !(c.Response().Status >= 200 && c.Response().Status <= 299) {
				logrus.
					WithField("http.header", header).
					WithField("method", c.Request().Method).
					WithField("traceID", traceID).
					WithField("path", c.Request().URL.Path).
					WithField("query", c.Request().URL.RawQuery).
					WithField("request.body", reqBodyBuf.String()).
					WithField("response.body", respBody).
					WithField("error", err).
					WithField("duration", duration.String()).
					Errorf("error request for client")

				return err
			}

			logrus.
				WithField("http.header", header).
				WithField("method", c.Request().Method).
				WithField("traceID", traceID).
				WithField("path", c.Request().URL.Path).
				WithField("query", c.Request().URL.RawQuery).
				WithField("request.body", reqBodyBuf.String()).
				WithField("response.body", respBody).
				WithField("error", err).
				WithField("duration", duration.String()).
				Infof("response success for client")

			return nil
		}
	}
}

type bodyCaptureWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
