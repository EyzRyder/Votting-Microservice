package controllers

import (
    "net/http"
    "html/template"
)

func Home_Controller (w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("view/index.html"))
    tmpl.Execute(w,nil)
}
