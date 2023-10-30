package template

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Template struct {
	Engine     string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
}

type Data struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	Secure          bool
	ServerName      string
}

func (t *Template) Render(w http.ResponseWriter, _ *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(t.Engine) {
	case "go":
		err := t.GoTemplate(w, view, data)
		if err != nil {
			return err
		}
	case "jet":
		err := t.JetTemplate(w, view, variables, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// GoTemplate renders a page using the Go template
func (t *Template) GoTemplate(w http.ResponseWriter, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/templates/%s.tmpl", t.RootPath, view))
	if err != nil {
		return err
	}

	td := &Data{}
	if data != nil {
		td = data.(*Data)
	}

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

// JetTemplate renders a page using the Jet templating engine
func (t *Template) JetTemplate(w http.ResponseWriter, view string, variables, data interface{}) error {
	// Jet template is not actively maintained. Support an alternative and make it default
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &Data{}
	if data != nil {
		td = data.(*Data)
	}

	tmpl, err := t.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = tmpl.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
