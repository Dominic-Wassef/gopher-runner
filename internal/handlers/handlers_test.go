package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var loginTests = []struct {
	name                 string
	url                  string
	method               string
	postedData           url.Values
	expectedResponseCode int
}{
	{
		name:                 "login-screen",
		url:                  "/",
		method:               "GET",
		expectedResponseCode: http.StatusOK,
	},
	{
		name:   "login-screen-post",
		url:    "/",
		method: "POST",
		postedData: url.Values{
			"email":    {"me@here.com"},
			"password": {"password"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
}

func TestLoginScreen(t *testing.T) {
	for _, e := range loginTests {
		if e.method == "GET" {
			req, _ := http.NewRequest(e.method, e.url, nil)

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Repo.LoginScreen)

			handler.ServeHTTP(rr, req)
			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		} else {
			req, _ := http.NewRequest(e.method, e.url, strings.NewReader(e.postedData.Encode()))

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(Repo.Login)

			handler.ServeHTTP(rr, req)
			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		}
	}
}
