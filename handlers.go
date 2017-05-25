package main

import "net/http"

import "io/ioutil"
import "encoding/json"

func (a *App) HandleNewScore(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var s score
	err = json.Unmarshal(b, &s)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = s.create(a.DB)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	j, _ := json.Marshal(s)
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func (a *App) HandleListScores(w http.ResponseWriter, r *http.Request) {
	scores, err := getScores(a.DB, 10, 0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	j, _ := json.Marshal(scores)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
