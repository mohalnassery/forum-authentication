package routes

import (
	"forum/models"
	"net/http"
	"strings"
	"text/template"
)

func HandleGet(w http.ResponseWriter, r *http.Request) {
	var err error
	var tmpl *template.Template
	path, _, _ := strings.Cut(r.URL.Path[1:], "/")
	switch path {
	case "register":
		tmpl, err = template.ParseFiles("../client/registerPage.html")
	case "login":
		tmpl, err = template.ParseFiles("../client/loginPage.html")
	case "":
		tmpl, err = template.ParseFiles("../client/index.html")
	case "create":
		tmpl, err = template.ParseFiles("../client/create-post.html")
	case "post-details":
		tmpl, err = template.ParseFiles("../client/postDetails.html")
	case "404":
		tmpl, err = template.ParseFiles("../client/errors/404.html")
		w.WriteHeader(http.StatusNotFound)
	case "400":
		tmpl, err = template.ParseFiles("../client/errors/400.html")
		w.WriteHeader(http.StatusBadRequest)
	case "500":
		tmpl, err = template.ParseFiles("../client/errors/500.html")
		w.WriteHeader(http.StatusInternalServerError)
	case "405":
		tmpl, err = template.ParseFiles("../client/errors/405.html")
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		tmpl, err = template.ParseFiles("../client/errors/404.html")
		w.WriteHeader(http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func CreatePostTemplate(w http.ResponseWriter, formdata models.PostCreate) {
	tmpl, err := template.ParseFiles("../client/create-post.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, formdata); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
