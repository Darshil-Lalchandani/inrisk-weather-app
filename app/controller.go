package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Controller struct {
	manager Manager
}

func NewController(m Manager) Controller {
	return Controller{
		manager: m,
	}
}

func (c Controller) MountRoutes(r chi.Router) {
	r.Post("/store-weather-data", c.FetchWeatherData)
	r.Get("/list-weather-files", c.ListWeatherFiles)
	r.Get("/weather-file-content/{fileName}", c.ListWeatherFileContent)
}

func (c Controller) FetchWeatherData(w http.ResponseWriter, r *http.Request) {
	var input FetchWeatherDataRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		c.ErrorWith(w, http.StatusBadRequest, err)
	}

	err = c.manager.FetchWeatherData(input)
	if err != nil {
		c.ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	res := FetchWeatherDataResponse{
		Message: "Success",
	}
	c.RespondWith(w, http.StatusOK, res)
}

func (c Controller) ListWeatherFiles(w http.ResponseWriter, r *http.Request) {
	files, err := c.manager.ListFiles()
	if err != nil {
		c.ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	c.RespondWith(w, http.StatusOK, files)

}

func (c Controller) ListWeatherFileContent(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")

	res, err := c.manager.FetchWeatherFileContent(fileName)
	if err != nil {
		c.ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	c.RespondWith(w, http.StatusOK, res)
}

func (c Controller) ErrorWith(w http.ResponseWriter, statusCode int, err error) {
	errMsg := struct {
		Error string `json:"message"`
	}{
		Error: err.Error(),
	}

	msg, _ := json.Marshal(errMsg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(msg)
}

func (c Controller) RespondWith(w http.ResponseWriter, statusCode int, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		c.ErrorWith(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
