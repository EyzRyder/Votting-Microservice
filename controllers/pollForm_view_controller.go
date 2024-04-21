package controllers

import (
    "net/http"
		"html/template"
)

func PollForm_Controller(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("view/PollForm.html"))
    tmpl.Execute(w,nil)
}
