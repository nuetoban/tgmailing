package view

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	s StatisticsProvider
}

func NewJSON(s StatisticsProvider) *JSON {
	return &JSON{s: s}
}

func (j *JSON) StatisticsHTTP(w http.ResponseWriter, req *http.Request) {
	s := j.s.Statistics()

	w.Header().Add("Content-Type", "application/json")

	content, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(Error{Error: err}.JSON())
		return
	}

	w.Write(content)
}
