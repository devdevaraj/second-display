package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"vdisplay/internal/service"

	"github.com/pion/webrtc/v4"
)

type Handler struct {
	Manager *service.SessionManager
}

// corsMiddleware adds basic CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// SetupRoutes registers the HTTP routes
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/sessions", corsMiddleware(h.handleSessions))
	mux.HandleFunc("/api/sessions/", corsMiddleware(h.handleSessionAction))
}

func (h *Handler) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sessions := h.Manager.ListSessions()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sessions)
		return
	}

	if r.Method == "POST" {
		var req struct {
			Resolution string `json:"resolution"`
			FPS        int    `json:"fps"`
			Bitrate    int    `json:"bitrate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sess, err := h.Manager.CreateSession(req.Resolution, req.FPS, req.Bitrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sess)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handler) handleSessionAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/sessions/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	id := parts[0]

	if r.Method == "DELETE" && len(parts) == 1 {
		err := h.Manager.DeleteSession(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "POST" && len(parts) == 2 {
		action := parts[1]
		ctx := r.Context()

		if action == "start" {
			if err := h.Manager.StartSession(ctx, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		} else if action == "stop" {
			if err := h.Manager.StopSession(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		} else if action == "webrtc" {
			var offer webrtc.SessionDescription
			if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			webrtcHandler, err := h.Manager.GetWebRTC(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			answer, err := webrtcHandler.HandleOffer(offer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(answer)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
