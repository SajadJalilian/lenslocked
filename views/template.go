package views

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func Parse(filePath string) (Template, error) {
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, filePath string, date interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, date)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
