// Package api implements the main clickhouse auth api
package api

import (
	"encoding/json"
	"net/http"

	"github.com/13excite/clickhouse-proxy-auth/pkg/version"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// handler implements the main api logic
type handler struct {
	hostToCluster    map[string]string
	aclClustersRules map[string][]string
	logger           *zap.SugaredLogger
}

// AuthResponse is the response for a authClickhouse request
type AuthResponse struct {
	Status string `json:"status"`
}

// HandlerOpt is an option func for the handler
type HandlerOpt func(*handler)

// NewHandler creates a new Clickhouse Auth API
func NewHandler(aclClusterRules map[string][]string, hostToCluster map[string]string) http.Handler {
	h := &handler{
		logger:           zap.S().With("package", "autocomplete"),
		hostToCluster:    hostToCluster,
		aclClustersRules: aclClusterRules,
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/auth", h.authClickhouse)
		r.Get("/version", h.serveVersion)
	})
	return router
}

// authClickhouse is the main handler for the authClickhouse request
func (h *handler) authClickhouse(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.Header.Get("X-Remote-IP")
	serverName := r.Header.Get("X-Server")
	if serverName == "" {
		h.logger.Warn("header X-Server not found.", "headers: ", r.Header)
		respondWithError(w, http.StatusForbidden, "header X-Server not found")
		return
	}
	clusterName, clusterNameOk := h.hostToCluster[serverName]
	if !clusterNameOk {
		h.logger.Warn("server not found in config ", "server ", serverName)
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	currentClusterRules := h.aclClustersRules[clusterName]
	allowAccess, err := checkIPInSubnet(remoteIP, currentClusterRules)
	if err != nil {
		h.logger.Warn("Could not parsing subnet ", "Error: ", err)
		respondWithError(w, http.StatusForbidden, "Invalid subnets")
		return
	}
	if !allowAccess {
		h.logger.Infow("subnets doesn't contains x-real-ip", "x-real-ip", remoteIP)
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}
	h.logger.Infow("Allow access", "x-real-ip", remoteIP, " to cluster: ", clusterName)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

// serveHTTP implements the http.Handler interface
// TODO: move to the separate pkg
func (h *handler) serveVersion(w http.ResponseWriter, _ *http.Request) {
	buildInfo := version.Build
	respondWithJSON(w, http.StatusOK, buildInfo)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
