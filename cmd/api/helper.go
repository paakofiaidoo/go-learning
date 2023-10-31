package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelop map[string]any

func returnJSON(w http.ResponseWriter, data envelop) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r http.Request, dst any) error {
	maxBytes := 1_048 + 576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only contain a single json object")
	}

	return nil
}
