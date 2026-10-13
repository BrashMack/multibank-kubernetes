package openapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// AuthTransport - интерсептор валидирует токен для запрашиваемого банка (по hostname запроса),
// обновляет его в случае истекания и добавляет токен в заголовок.
func (c *Client) authTransport(next http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		bank, ok := c.bankForHostname[r.URL.Hostname()]
		if !ok {
			return next.RoundTrip(r)
		}

		token, ok := c.tokenForBank[bank]
		if !ok {
			return next.RoundTrip(r)
		}

		if time.Now().After(token.ExpiredAt) {
			var baseUrl string

			if baseUrl, ok = c.baseUrlForBank[bank]; !ok {
				return next.RoundTrip(r)
			}

			accessToken, err := c.getToken(baseUrl)
			if err != nil {
				logrus.
					WithField("url", baseUrl).
					WithField("bank", bank).
					WithField("error", err).
					Error("c.getToken response error")

				return next.RoundTrip(r)
			}

			c.tokenForBank[bank] = accessToken
		}

		r.Header.Set("Authorization", "Bearer "+c.tokenForBank[bank].AccessToken)

		return next.RoundTrip(r)
	})
}

func (c *Client) logTransport(next http.RoundTripper, hiddenHeaders []string) http.RoundTripper {
	var hidden = append(slices.Clone(hiddenHeaders), "authorization")

	return roundTripperFunc(func(parent *http.Request) (*http.Response, error) {
		traceID, ok := parent.Context().Value("traceID").(string)
		if !ok || traceID == "" {
			traceID = uuid.New().String()
		}

		hdr := make(http.Header, len(parent.Header))
		for key, value := range parent.Header {
			cp := slices.Clone(value)
			hdr[strings.ToLower(key)] = cp
		}

		for _, h := range hidden {
			if hdr.Get(strings.ToLower(h)) != "" {
				hdr[strings.ToLower(h)] = []string{"*"}
			}
		}

		logrus.
			WithField("traceID", traceID).
			WithField("method", parent.Method).
			WithField("url", parent.URL.String()).
			WithField("header", hdr).
			Info("outing request to openapi")

		ts := time.Now()
		resp, err := next.RoundTrip(parent)
		dur := time.Since(ts)

		if err != nil {
			logrus.
				WithField("traceID", traceID).
				WithField("method", parent.Method).
				WithField("url", parent.URL.String()).
				WithField("header", hdr).
				WithField("error", err.Error()).
				WithField("duration", dur.String()).
				Error("response error from openapi")

			return nil, err
		}

		var respBody []byte

		if resp.Body != nil {
			respBody, _ = io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody)) // восстанавливаем
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			logrus.WithField("traceID", traceID).
				WithField("method", parent.Method).
				WithField("url", parent.URL.String()).
				WithField("header", hdr).
				WithField("statusCode", resp.StatusCode).
				WithField("duration", dur.String()).
				WithField("responseBody", string(respBody)).
				Info("response success from openapi")
		} else {
			logrus.WithField("traceID", traceID).
				WithField("method", parent.Method).
				WithField("url", parent.URL.String()).
				WithField("header", hdr).
				WithField("statusCode", resp.StatusCode).
				WithField("duration", dur.String()).
				WithField("responseBody", string(respBody)).
				Error("response error from openapi")
		}

		return resp, nil
	})
}

func (c *Client) retryTransport(next http.RoundTripper, codes []int, period, deadline time.Duration) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		ctx, cancel := context.WithTimeout(r.Context(), deadline)
		defer cancel()

		req := r.Clone(ctx)

		var invoke = func() (*http.Response, error) {
			res, err := next.RoundTrip(req)
			if err != nil {
				return nil, fmt.Errorf("next.RoundTrip: %w", err)
			}

			if res.StatusCode >= 200 && res.StatusCode <= 299 {
				return res, nil
			}

			if !isRetryableHTTP(res.StatusCode, codes) {
				return res, nil
			}

			_ = res.Body.Close()

			return nil, fmt.Errorf("http status code %d", res.StatusCode)
		}

		resp, err := invoke()

		if !retryAllowedHTTP(err, period, deadline) {
			return resp, err
		}

		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				if err == nil {
					err = ctx.Err()
				}

				return nil, err
			case <-ticker.C:
				resp, err = invoke()

				if !retryAllowedHTTP(err, period, deadline) {
					return resp, err
				}
			}
		}
	})
}

func retryAllowedHTTP(err error, period, deadline time.Duration) bool {
	if period <= 0 || deadline <= 0 {
		return false
	}

	return err != nil
}

func isRetryableHTTP(code int, codes []int) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}

	return false
}
