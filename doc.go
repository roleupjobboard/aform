/*
Package aform provides functionalities for working effectively with forms.

Introduction

This package provides simple tools to create forms, to render
them as HTML and to do forms' input validations.

Let's suppose you want to create a form to get the user's name.
You'll need something like this in your template.

	<form action="/your-name" method="post">
		<label for="your_name">Your name: </label>
		<input id="your_name" type="text" name="your_name" value="{{ .CurrentName }}">
		<input type="submit" value="OK">
	</form>

This HTML code tells the browser to return the form data to the URL /your-name,
using the POST method. It will display a text field labeled "Your name:" and a
button "OK". If the template context contains the CurrentName variable, it will
be used to pre-fill the your_name field.

You'll need a handler to render the template containing the HTML form
and providing the CurrentName value.

When the form is submitted, the POST request which is sent to the server
will contain the form data. So, you'll also need a handler for the POST request
to the path /your-name. This handler will have to find the appropriate
key/values pairs in the request and to process them.

Form Creation

We already know what we want our HTML form to look like. Our starting point to
create it with aform is this:
	nameForm := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("Your name"))),
	))
This creates a Form with a single CharField. We used DefaultCharField that
creates a CharField with reasonable defaults like a maximum length of 256.
It means that the browser should prevent the user from entering more
characters than that. It also means that when your application HTTP handler
receives the form back from the browser, aform will validate the length
of the data.

A Form has a Form.IsValid function, which runs validation for all its fields.
When this function is called, if all fields contain valid data, it will return
true and make the form's data accessible via Form.CleanedData function.

The whole nameForm, when rendered for the first time, will look like:
	<div><label for="id_your_name">Your name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" required></div>
Note that it does not include the <form> tags, or a submit button. We’ll have
to provide those ourselves in the template.

The HTTP handler

Form data sent back to the server is processed by a handler. To handle the form
we need to create the form in the handler for the path where we want it to be
published:

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

If the handler is triggered by a GET request, it will create an empty form
and place it in the template context to be rendered. This is what we can expect
to happen the first time we visit the URL.

If the form is submitted using a POST request, the handler will populate the
form with data from the request:
	nameForm.BindRequest(req)
This is called “binding data to the form” (it is now a bound form).

We call the Form.IsValid method; if it’s not True, we go back to the template
with the form. This time the form is no longer empty (unbound) so the HTML form
will be populated with the data previously submitted, where it can be edited
and corrected as required. The HTML form will also display errors found by the
validation.

If Form.IsValid is True, we’ll now be able to find all the validated form data
in its CleanedData returned by Form.CleanedData. We can use this data to update
the database or do other processing before sending an HTTP redirect to the
browser telling it where to go next.

The template

We don’t need to do much in our template:
*/
//	<form action="/your-name" method="post">
//		{{ .form.AsDiv }}
//		{{/* Add CSRF */}}
//		<button type="submit">OK</button>
//	</form>
/*

All the form’s fields and their attributes will be unpacked into HTML markup from that:
	{{ .form.AsDiv }}

Bound and unbound forms

A Form is either bound to a set of data, or unbound.
If it’s bound to a set of data, it’s capable of validating that data and
rendering the form as HTML with the data displayed in the HTML.
If it’s unbound, it cannot do validation (because there’s no data to
validate), but it can still render the blank form as HTML.

To bind a Form, we use the method Form.BindRequest. Method Form.BindData can
be used to when the application needs to modify the default behavior.

*/
package aform
