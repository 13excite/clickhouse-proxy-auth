package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedError struct {
	ErrMsg string `json:"error"`
}

func TestNewHandler(t *testing.T) {
	//validStatusCodes := []int{http.StatusOK, http.StatusForbidden}

	tests := []struct {
		name             string
		aclClusterRules  map[string][]string
		hostToCluster    map[string]string
		method           string
		path             string
		requestXServer   string
		requstXRemoteIP  string
		wantErr          bool
		expectedErrorMsg string
		wantStatusCode   int
	}{
		{
			name: "requst from correct subnet and to server",
			aclClusterRules: map[string][]string{
				"data_science": {"10.10.10.0/24", "100.10.20.0/24"}},
			hostToCluster:   map[string]string{"clickhouse-1": "data_science"},
			method:          http.MethodGet,
			path:            "/auth",
			requestXServer:  "clickhouse-1",
			requstXRemoteIP: "10.10.10.10",
			wantErr:         false,
			wantStatusCode:  http.StatusOK,
		},
		{
			name: "requst from incorrect subnet and to server",
			aclClusterRules: map[string][]string{
				"data_science": {"10.10.10.0/24", "100.10.20.0/24"}},
			hostToCluster:   map[string]string{"clickhouse-1": "data_science"},
			method:          http.MethodGet,
			path:            "/auth",
			requestXServer:  "clickhouse-1",
			requstXRemoteIP: "10.30.10.10",
			wantErr:         false,
			wantStatusCode:  http.StatusForbidden,
		},
		{
			name: "requst from correct subnet and to incorrect server",
			aclClusterRules: map[string][]string{
				"data_science": {"10.10.10.0/24", "100.10.20.0/24"},
				"backend":      {"10.40.40.0/24"},
			},
			hostToCluster: map[string]string{
				"clickhouse-1": "data_science",
				"clickhouse-2": "backend",
			},
			method:          http.MethodGet,
			path:            "/auth",
			requestXServer:  "clickhouse-2",
			requstXRemoteIP: "10.10.10.15",
			wantErr:         false,
			wantStatusCode:  http.StatusForbidden,
		},
		{
			name: "error invalid subnets expected",
			aclClusterRules: map[string][]string{
				"data_science": {"10.300.10.0/24", "100.300.20.0/24"}},
			hostToCluster:    map[string]string{"clickhouse-1": "data_science"},
			method:           http.MethodGet,
			path:             "/auth",
			requestXServer:   "clickhouse-1",
			requstXRemoteIP:  "300.290.10.10",
			wantErr:          true,
			expectedErrorMsg: "Invalid subnets",
			wantStatusCode:   http.StatusForbidden,
		},
		{
			name: "error x-server header not found expected",
			aclClusterRules: map[string][]string{
				"data_science": {"10.10.10.0/24", "100.10.20.0/24"}},
			hostToCluster:    map[string]string{"clickhouse-1": "data_science"},
			method:           http.MethodGet,
			path:             "/auth",
			requestXServer:   "",
			requstXRemoteIP:  "10.10.10.10",
			wantErr:          true,
			expectedErrorMsg: "header X-Server not found",
			wantStatusCode:   http.StatusForbidden,
		},
		// add case >>>>header X-Server not found<<<<<
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// define handler
			h := NewHandler(tt.aclClusterRules, tt.hostToCluster)
			// start test server
			testServer := httptest.NewServer(h)
			defer testServer.Close()
			// prepare request
			req, err := http.NewRequest(tt.method, testServer.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("failed to set up test request: %s", err)
			}
			// set headers
			req.Header.Set("X-Remote-IP", tt.requstXRemoteIP)
			req.Header.Set("X-Server", tt.requestXServer)

			res, err := testServer.Client().Do(req)
			if err != nil {
				t.Fatalf("failed to perform testServer request: %s", err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatusCode, res.StatusCode, "status code should be equal")

			if tt.wantErr {
				var jsonError expectedError
				if err := json.NewDecoder(res.Body).Decode(&jsonError); err != nil {
					t.Errorf("returned error can't be decoded onto JSONError: %s", err)
				}
				assert.Equal(t, tt.expectedErrorMsg, jsonError.ErrMsg, "error message should be equal")
				return
			}
		})
	}
}
