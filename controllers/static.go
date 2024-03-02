package controllers

import (
	"html/template"
	"net/http"

	"github.com/sajadjalilian/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "First Q",
			Answer:   "First Q answer",
		},
		{
			Question: "Second Q",
			Answer:   "Second Q answer",
		},
		{
			Question: "Third Q",
			Answer:   "Third Q answer",
		},
		{
			Question: "What is your email?",
			Answer:   `<a href="mailto:sajadjalilian88@gmail.com">sajadjalilian88@gmail.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
