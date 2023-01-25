package handlers

import (
	"encoding/json"
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

func TestDBRepo_PusherAuth(t *testing.T) {
	postedData := url.Values{
		"socket_id":    {"1359961149.1836558958"},
		"channel_name": {"private-channel-1"},
	}

	req, _ := http.NewRequest("POST", "/pusher/auth", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PusherAuth)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected response 200 but got %d", rr.Code)
	}

	type pusherResp struct {
		Auth string `json:"auth"`
	}

	var p pusherResp

	err := json.NewDecoder(rr.Body).Decode(&p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Auth) == 0 {
		t.Error("Empty JSON response!")
	}
}
