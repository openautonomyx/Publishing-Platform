//go:build draft

// Excluded from the default build: not routed from main.go and depends on code
// that does not exist yet. Build with -tags draft to work on it (see README).

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var deploymentTracker *DeploymentTracker

func init() {
	deploymentTracker = NewDeploymentTracker()
}

// RegisterDeploymentTrackingRoutes registers deployment tracking endpoints
func RegisterDeploymentTrackingRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/deployments/start", StartDeploymentHandler).Methods("POST")
	router.HandleFunc("/api/v1/deployments/{id}", GetDeploymentHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/{id}/progress", GetDeploymentProgressHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments", ListDeploymentsHandler).Methods("GET")
	router.HandleFunc("/api/v1/deployments/{id}/stream", StreamDeploymentEventsHandler).Methods("GET")
}

// StartDeploymentHandler starts a new deployment
func StartDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VPSHost string `json:"vps_host"`
		Branch  string `json:"branch"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.VPSHost == "" {
		req.VPSHost = "agennext.com"
	}
	if req.Branch == "" {
		req.Branch = "main"
	}

	deployment, err := deploymentTracker.StartDeployment(r.Context(), req.VPSHost, req.Branch)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start deployment")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"deployment_id": deployment.ID.String(),
		"status":        string(deployment.Status),
		"progress":      deployment.Progress,
		"total_steps":   deployment.TotalSteps,
		"started_at":    deployment.StartedAt,
	})
}

// GetDeploymentHandler retrieves deployment status
func GetDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	id, err := parseUUID(deploymentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deployment ID")
		return
	}

	deployment, err := deploymentTracker.GetDeployment(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Deployment not found")
		return
	}

	response := map[string]interface{}{
		"id":               deployment.ID.String(),
		"vps_host":         deployment.VPSHost,
		"branch":           deployment.Branch,
		"status":           string(deployment.Status),
		"progress":         deployment.Progress,
		"total_steps":      deployment.TotalSteps,
		"started_at":       deployment.StartedAt,
		"completed_at":     deployment.CompletedAt,
		"current_phase":    string(deployment.CurrentPhase),
		"failed_phase":     deployment.FailedPhase,
		"error_message":    deployment.ErrorMessage,
		"events_count":     len(deployment.Events),
	}

	if len(deployment.Events) > 0 {
		response["events"] = deployment.Events
	}

	respondJSON(w, http.StatusOK, response)
}

// GetDeploymentProgressHandler returns real-time progress
func GetDeploymentProgressHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	id, err := parseUUID(deploymentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deployment ID")
		return
	}

	progress := deploymentTracker.GetDeploymentProgress(id)
	if progress == nil {
		respondError(w, http.StatusNotFound, "Deployment not found")
		return
	}

	respondJSON(w, http.StatusOK, progress)
}

// ListDeploymentsHandler lists recent deployments
func ListDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		var parsedLimit int
		if _, err := scanf(l, "%d", &parsedLimit); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	deployments, err := deploymentTracker.ListDeployments(r.Context(), limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list deployments")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deployments": deployments,
		"count":       len(deployments),
	})
}

// StreamDeploymentEventsHandler streams deployment events (Server-Sent Events)
func StreamDeploymentEventsHandler(w http.ResponseWriter, r *http.Request) {
	deploymentID := mux.Vars(r)["id"]

	id, err := parseUUID(deploymentID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deployment ID")
		return
	}

	// Check if deployment exists
	deployment, err := deploymentTracker.GetDeployment(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Deployment not found")
		return
	}

	// Setup SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	// Stream current events
	for _, event := range deployment.Events {
		eventData, _ := json.Marshal(event)
		w.Write([]byte("data: "))
		w.Write(eventData)
		w.Write([]byte("\n\n"))
		flusher.Flush()
	}

	// Keep connection open and stream new events
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	lastEventCount := len(deployment.Events)

	for {
		select {
		case <-r.Context().Done():
			return

		case <-ticker.C:
			// Refresh deployment to get new events
			currentDeployment, err := deploymentTracker.GetDeployment(r.Context(), id)
			if err != nil {
				continue
			}

			// Send any new events
			if len(currentDeployment.Events) > lastEventCount {
				for _, event := range currentDeployment.Events[lastEventCount:] {
					eventData, _ := json.Marshal(event)
					w.Write([]byte("data: "))
					w.Write(eventData)
					w.Write([]byte("\n\n"))
					flusher.Flush()
				}
				lastEventCount = len(currentDeployment.Events)
			}

			// Check if deployment is complete
			if currentDeployment.Status != "IN_PROGRESS" {
				// Send final status
				finalData, _ := json.Marshal(map[string]interface{}{
					"type":   "complete",
					"status": string(currentDeployment.Status),
				})
				w.Write([]byte("data: "))
				w.Write(finalData)
				w.Write([]byte("\n\n"))
				flusher.Flush()
				return
			}
		}
	}
}

// Helper functions
func parseUUID(s string) (UUID, error) {
	return uuid.Parse(s)
}

func scanf(input string, format string, args ...interface{}) (int, error) {
	// Simple implementation for parsing integers
	_, err := fmt.Sscanf(input, format, args...)
	return 0, err
}
