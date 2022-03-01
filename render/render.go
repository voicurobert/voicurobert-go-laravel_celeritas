package render

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (r *Render) defaultData(td *TemplateData, req *http.Request) {
	td.Secure = r.Secure
	td.ServerName = r.ServerName
	td.Port = r.Port
	if r.Session.Exists(req.Context(), "userID") {
		td.IsAuthenticated = true
	}
}

func (r *Render) Page(w http.ResponseWriter, req *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(r.Renderer) {
	case "go":
		return r.GoPage(w, req, view, data)
	case "jet":
		return r.JetPage(w, req, view, variables, data)
	default:

	}
	return errors.New("no rendering engine specified")
}

// GoPage renders a standard go template
func (r *Render) GoPage(w http.ResponseWriter, req *http.Request, templateName string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", r.RootPath, templateName))
	if err != nil {
		return err
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}
	return nil
}

// JetPage renders a template using jet
func (r *Render) JetPage(w http.ResponseWriter, req *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	r.defaultData(td, req)

	t, err := r.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
