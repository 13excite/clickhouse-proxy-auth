// Package api implements the main clickhouse auth api
package api

import (
	"encoding/json"
	"net"
	"net/http"

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
func NewHandler(opts ...HandlerOpt) http.Handler {
	h := &handler{}
	for _, o := range opts {
		o(h)
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/auth", h.authClickhouse)
	})
	return router
}

// TODO: add logic
func (h *handler) authClickhouse(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.Header.Get("X-Remote-IP")
	serverName := r.Header.Get("X-Server")
	if serverName == "" {
		h.logger.Warn("header X-Server not found.", "Headers: ", r.Header)
		respondWithError(w, http.StatusForbidden, "header X-Server not found")
		return
	}
	clusterName, clusterNameOk := h.hostToCluster[serverName]
	if !clusterNameOk {
		h.logger.Warn("server not found in config ", "Server", serverName)
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	currentClusterRules := h.aclClustersRules[clusterName]
	allowAccess, err := checkIpInSubnet(remoteIp, currentClusterRules)
	if err != nil {
		h.logger.Warn("Could not parsing subnet ", "Error: ", err)
		respondWithError(w, http.StatusForbidden, "Invalid subnets")
		return
	}
	if !allowAccess {
		h.logger.Infow("subnets doesn't contains x-real-ip", "x-real-ip", remoteIp)
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}
	h.logger.Infow("Allow access", "x-real-ip", remoteIp, " to cluster: ", clusterName)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})

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

// checkIpInSubnet checks contains ip in subnet
func checkIpInSubnet(ipAddr string, subnets []string) (bool, error) {
	// iterate by subnets array and check
	// does subnet contain addr or not
	for _, subnet := range subnets {
		_, subnetParse, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}
		ipAddrParse := net.ParseIP(ipAddr)
		if subnetParse.Contains(ipAddrParse) {
			return true, nil
		} // end if contains
	} // end for

	return false, nil
}
