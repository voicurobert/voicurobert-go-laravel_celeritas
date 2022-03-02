package celeritas

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

func (c *Celeritas) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, filename string) error {
	fp := path.Join(pathToFile, filename)
	fileToServe := filepath.Clean(fp)
	//w.Header().Set("Content-Type", fmt.Sprintf("attachment; file=\"%s\"", filename))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", filename))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (c *Celeritas) Error404(w http.ResponseWriter) {
	c.ErrorStatus(w, http.StatusNotFound)
}

func (c *Celeritas) Error500(w http.ResponseWriter) {
	c.ErrorStatus(w, http.StatusInternalServerError)
}

func (c *Celeritas) ErrorUnauthorized(w http.ResponseWriter) {
	c.ErrorStatus(w, http.StatusUnauthorized)
}

func (c *Celeritas) ErrorForbidden(w http.ResponseWriter) {
	c.ErrorStatus(w, http.StatusForbidden)
}

func (c *Celeritas) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
