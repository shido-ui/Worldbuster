package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"net/http"
)

type jobHandler struct{ s *job.Service }

func newJobHandler(s *job.Service) *jobHandler                    { return &jobHandler{s: s} }
func (h *jobHandler) list(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, h.s.List()) }
func (h *jobHandler) employ(w http.ResponseWriter, r *http.Request) {
	var q struct {
		CharacterID string `json:"characterId"`
		JobID       string `json:"jobId"`
		PositionID  string `json:"positionId"`
		Education   int    `json:"education"`
		Stat        int    `json:"stat"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	if err := h.s.Employ(q.CharacterID, q.JobID, q.PositionID, q.Education, q.Stat); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	e, _ := h.s.GetEmployment(q.CharacterID)
	writeJSON(w, 200, e)
}
