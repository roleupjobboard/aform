package aform_test

import (
	"fmt"
	"github.com/roleupjobboard/aform"
	"html/template"
	"net/http"
)

func Example_yourName() {
	http.HandleFunc("/your-name", nameHandler)
	http.HandleFunc("/congratulations", congratulationsHandler)
	panic(http.ListenAndServe(":8080", nil))
}

func nameHandler(w http.ResponseWriter, req *http.Request) {
	nameForm := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("Your name"))),
	))
	// if this is a POST request we need to process the form data
	if req.Method == "POST" {
		// Populate the form with data from the request:
		nameForm.BindRequest(req)
		// check whether it's valid:
		if nameForm.IsValid() {
			// process the data with Form.CleanedData
			// ...
			// redirect to a new URL:
			http.Redirect(w, req, "/congratulations", http.StatusFound)
		}
	}
	_ = tmpl.ExecuteTemplate(w, "name", map[string]any{"form": nameForm})
}

var tmpl = template.Must(template.New("name").Parse(tmplBody))

const tmplBody = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Name Form</title>
	</head>
	<body>
		<form action="/your-name" method="post">
			{{ .form.AsDiv }}
			{{/* Add CSRF */}}
			<button type="submit">OK</button>
		</form>
	</body>
</html>`

func congratulationsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Congratulations!")
}
