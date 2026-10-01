package api

import (
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
	"strconv"
)

type NewsAPI struct{ Repo store.NewsRepository }

func (a *NewsAPI) Recent(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if n, e := strconv.Atoi(r.URL.Query().Get("limit")); e == nil && n > 0 {
		limit = n
	}
	x, err := a.Repo.Recent(r.Context(), limit)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "news unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"news": x})
}
