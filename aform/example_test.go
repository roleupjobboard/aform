package aform_test

import (
	"fmt"
)

func ExampleNew_yourName() {
	nameForm := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("Your name"))),
	))
	fmt.Println(nameForm.AsDiv())
	// Output:
	// <div><label for="id_your_name">Your name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" required></div>
}

func ExampleNew_login() {
	f := aform.Must(aform.New(
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Password", aform.WithWidget(aform.PasswordInput)))),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Remember me"))),
	))
	fmt.Println(f.AsDiv())
	// Output:
	// <div><label for="id_email">Email</label><input type="email" name="email" id="id_email" maxlength="254" required></div>
	// <div><label for="id_password">Password</label><input type="password" name="password" id="id_password" maxlength="256" required></div>
	// <div><label for="id_remember_me">Remember me</label><input type="checkbox" name="remember_me" id="id_remember_me" required></div>
}

func ExampleNew_notRequiredField() {
	// Create a form with two not required fields
	nicknameFld, _ := aform.NewCharField("Nickname", "", "nickname_default", 0, 0, aform.IsNotRequired())
	emailFld, _ := aform.NewEmailField("Email", "", "catchall@domain.com", 0, 0, aform.IsNotRequired())
	f, err := aform.New(
		aform.WithCharField(nicknameFld),
		aform.WithEmailField(emailFld),
	)
	if err != nil {
		panic("invalid form")
	}
	// Simulate an empty form post. Bind nil
	f.BindData(nil)
	// IsValid is true because fields are not required
	if !f.IsValid() {
		panic("form must be valid")
	}
	fmt.Println(f.CleanedData().Get("nickname"))
	fmt.Println(f.CleanedData().Get("email"))
	// Output:
	// nickname_default
	// catchall@domain.com
}

func ExampleNew_profile() {
	f, err := aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("Username"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Avatar", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
			{Value: "apple", Label: "Apple"},
			{Value: "orange", Label: "Orange"},
			{Value: "strawberry", Label: "Strawberry"},
		})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField(
			"Color",
			aform.WithGroupedChoiceOptions("Basic", []aform.ChoiceFieldOption{
				{Value: "red", Label: "Red"},
				{Value: "green", Label: "Green"},
				{Value: "blue", Label: "Blue"},
			}),
			aform.WithGroupedChoiceOptions("Fancy", []aform.ChoiceFieldOption{
				{Value: "orange", Label: "Orange"},
				{Value: "purple", Label: "Purple"},
				{Value: "pink", Label: "Pink"},
			}),
		)),
		),
	)
	if err != nil {
		panic("invalid form")
	}
	fmt.Println(f.AsDiv())
	// Output:
	// <div><label for="id_username">Username</label><input type="text" name="username" id="id_username" maxlength="256" required></div>
	// <div><label for="id_avatar">Avatar</label><select name="avatar" id="id_avatar" required>
	//   <option value="apple" id="id_avatar_0">Apple</option>
	//   <option value="orange" id="id_avatar_1">Orange</option>
	//   <option value="strawberry" id="id_avatar_2">Strawberry</option>
	// </select></div>
	// <div><label for="id_color">Color</label><select name="color" id="id_color" required>
	//   <optgroup label="Basic">
	//     <option value="red" id="id_color_0_0">Red</option>
	//     <option value="green" id="id_color_0_1">Green</option>
	//     <option value="blue" id="id_color_0_2">Blue</option>
	//   </optgroup>
	//   <optgroup label="Fancy">
	//     <option value="orange" id="id_color_1_0">Orange</option>
	//     <option value="purple" id="id_color_1_1">Purple</option>
	//     <option value="pink" id="id_color_1_2">Pink</option>
	//   </optgroup>
	// </select></div>
}

func ExampleForm_CleanedData_allValid() {
	f := aform.Must(aform.New(
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Password", aform.WithWidget(aform.PasswordInput)))),
	))
	f.BindData(map[string][]string{"email": {"john.doe@gmail.com"}, "password": {"Password123"}})
	f.IsValid()
	fmt.Printf("email: \"%s\"\n", f.CleanedData().Get("email"))
	fmt.Printf("password: \"%s\"\n", f.CleanedData().Get("password"))
	// Output:
	// email: "john.doe@gmail.com"
	// password: "Password123"
}

func ExampleForm_CleanedData_withInvalidField() {
	f := aform.Must(aform.New(
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Password", aform.WithWidget(aform.PasswordInput)))),
	))
	f.BindData(map[string][]string{"email": {"invalid email"}, "password": {"Password123"}})
	f.IsValid()
	fmt.Printf("email: \"%s\"\n", f.CleanedData().Get("email"))
	fmt.Printf("password: \"%s\"\n", f.CleanedData().Get("password"))
	fmt.Printf("email error: \"%s\"\n", f.Errors().Get("email"))
	fmt.Printf("password error: \"%s\"\n", f.Errors().Get("password"))
	// Output:
	// email: ""
	// password: "Password123"
	// email error: "Enter a valid email address"
	// password error: ""
}

func ExampleForm_SetCleanFunc() {
	f := aform.Must(aform.New(
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("email1", aform.WithLabel("Email")))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("email2", aform.WithLabel("Email verification")))),
	))
	f.SetCleanFunc(func(form *aform.Form) {
		if !form.IsValid() {
			return
		}
		data := form.CleanedData()
		if data.Get("email1") != data.Get("email2") {
			err := fmt.Errorf("Emails must be the same")
			_ = form.AddError("email1", err)
			_ = form.AddError("email2", err)
		}
	})
	f.BindData(map[string][]string{"email1": []string{"g@gmail.com"}, "email2": []string{"h@gmail.com"}})
	// Run the form validation
	f.IsValid()
	fmt.Println(f.AsDiv())
	// Output:
	// <div><label for="id_email1">Email</label>
	// <ul class="errorlist"><li id="err_0_id_email1">Emails must be the same</li></ul>
	// <input type="email" name="email1" value="g@gmail.com" id="id_email1" maxlength="254" aria-describedby="err_0_id_email1" aria-invalid="true" required></div>
	// <div><label for="id_email2">Email verification</label>
	// <ul class="errorlist"><li id="err_0_id_email2">Emails must be the same</li></ul>
	// <input type="email" name="email2" value="h@gmail.com" id="id_email2" maxlength="254" aria-describedby="err_0_id_email2" aria-invalid="true" required></div>
}

func ExampleWithLabelSuffix() {
	f := aform.Must(aform.New(
		aform.WithLabelSuffix(":"),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("email"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("password", aform.WithWidget(aform.PasswordInput)))),
	))
	fmt.Println(f.AsDiv())
	// Output:
	// <div><label for="id_email">email:</label><input type="email" name="email" id="id_email" maxlength="254" required></div>
	// <div><label for="id_password">password:</label><input type="password" name="password" id="id_password" maxlength="256" required></div>
}

func ExampleWithAttributes_classWithErrorClass() {
	f := aform.Must(aform.New(
		aform.WithRequiredCSSClass("required"),
		aform.WithErrorCSSClass("error"),
		aform.WithCharField(aform.Must(aform.DefaultCharField(
			"Name",
			aform.WithAttributes([]aform.Attributable{aform.Attr("class", "big blue")}),
		))),
	))
	// Bind empty data. As the field is required,
	// it will make IsValid == false and add error to the HTML generated.
	f.BindData(map[string][]string{})
	// Run the form validation
	f.IsValid()
	fmt.Println(f.AsDiv())
	// Output:
	// <div class="required error"><label class="required" for="id_name">Name</label>
	// <ul class="errorlist"><li id="err_0_id_name">This field is required</li></ul>
	// <input type="text" name="name" class="error big blue" id="id_name" maxlength="256" aria-describedby="err_0_id_name" aria-invalid="true" required></div>
}

func ExampleWithWidget_passwordInput() {
	fld := aform.Must(
		aform.DefaultCharField(
			"Password",
			aform.WithWidget(aform.PasswordInput),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_password">Password</label><input type="password" name="password" id="id_password" maxlength="256" required></div>
}

func ExampleWithWidget_textArea() {
	fld := aform.Must(
		aform.DefaultCharField(
			"Description",
			aform.WithWidget(aform.TextArea),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_description">Description</label><textarea name="description" id="id_description" maxlength="256" required>
	// </textarea></div>
}

func ExampleWithWidget_radioSelect() {
	fld := aform.Must(
		aform.DefaultChoiceField(
			"Color",
			aform.WithWidget(aform.RadioSelect),
			aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Red"},
				{Value: "green", Label: "Green"},
				{Value: "blue", Label: "Blue"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div>
	// <fieldset><legend for="id_color">Color</legend>
	// <div id="id_color">
	// <label for="id_color_0"><input type="radio" name="color" value="red" id="id_color_0">Red</label>
	// <label for="id_color_1"><input type="radio" name="color" value="green" id="id_color_1">Green</label>
	// <label for="id_color_2"><input type="radio" name="color" value="blue" id="id_color_2">Blue</label>
	// </div>
	// </fieldset>
	// </div>
}

func ExampleDefaultCharField_passwordWithPlaceholder() {
	fld := aform.Must(aform.DefaultCharField(
		"Password",
		aform.WithWidget(aform.PasswordInput),
		aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "****")})))
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_password">Password</label><input type="password" name="password" id="id_password" maxlength="256" placeholder="****" required></div>
}

func ExampleNewCharField_passwordWithPlaceholder() {
	fld := aform.Must(aform.NewCharField("Password", "", "", 0, 0,
		aform.WithWidget(aform.PasswordInput),
		aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "****")})))
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_password">Password</label><input type="password" name="password" id="id_password" placeholder="****" required></div>
}

func ExampleNewCharField_withInitialValue() {
	fld := aform.Must(aform.NewCharField("Your color", "white", "", 0, 0))
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_your_color">Your color</label><input type="text" name="your_color" value="white" id="id_your_color" required></div>
}

func ExampleDefaultChoiceField() {
	fld := aform.Must(
		aform.DefaultChoiceField("Color",
			aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_color">Color</label><select name="color" id="id_color" required>
	//   <option value="red" id="id_color_0">Rouge</option>
	//   <option value="green" id="id_color_1">Vert</option>
	//   <option value="blue" id="id_color_2">Bleu</option>
	// </select></div>
}

func ExampleDefaultChoiceField_withRadio() {
	fld := aform.Must(
		aform.DefaultChoiceField("Color",
			aform.WithWidget(aform.RadioSelect),
			aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div>
	// <fieldset><legend for="id_color">Color</legend>
	// <div id="id_color">
	// <label for="id_color_0"><input type="radio" name="color" value="red" id="id_color_0">Rouge</label>
	// <label for="id_color_1"><input type="radio" name="color" value="green" id="id_color_1">Vert</label>
	// <label for="id_color_2"><input type="radio" name="color" value="blue" id="id_color_2">Bleu</label>
	// </div>
	// </fieldset>
	// </div>
}

func ExampleNewChoiceField() {
	fld := aform.Must(
		aform.NewChoiceField("Color", "",
			aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "red", Label: "Rouge"}, {Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_color">Color</label><select name="color" id="id_color" required>
	//   <option value="red" id="id_color_0">Rouge</option>
	//   <option value="green" id="id_color_1">Vert</option>
	//   <option value="blue" id="id_color_2">Bleu</option>
	// </select></div>
}

func ExampleNewChoiceField_withRadio() {
	fld := aform.Must(
		aform.NewChoiceField("Color", "",
			aform.WithWidget(aform.RadioSelect),
			aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div>
	// <fieldset><legend for="id_color">Color</legend>
	// <div id="id_color">
	// <label for="id_color_0"><input type="radio" name="color" value="red" id="id_color_0">Rouge</label>
	// <label for="id_color_1"><input type="radio" name="color" value="green" id="id_color_1">Vert</label>
	// <label for="id_color_2"><input type="radio" name="color" value="blue" id="id_color_2">Bleu</label>
	// </div>
	// </fieldset>
	// </div>
}

func ExampleWithGroupedChoiceOptions() {
	fld := aform.Must(
		aform.DefaultChoiceField("Color",
			aform.WithGroupedChoiceOptions("RGB", []aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}),
		),
	)
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_color">Color</label><select name="color" id="id_color" required>
	//   <optgroup label="RGB">
	//     <option value="red" id="id_color_0_0">Rouge</option>
	//     <option value="green" id="id_color_0_1">Vert</option>
	//     <option value="blue" id="id_color_0_2">Bleu</option>
	//   </optgroup>
	// </select></div>
}

func ExampleField_SetLabel() {
	fld := aform.Must(aform.DefaultBooleanField("remember"))
	fld.SetLabel("Remember me")
	fmt.Println(fld.LabelTag())
	// Output:
	// <label for="id_remember">Remember me</label>
}

func ExampleField_SetLabelSuffix() {
	fld := aform.Must(aform.DefaultBooleanField("Remember me"))
	fld.SetLabelSuffix(":")
	fmt.Println(fld.LabelTag())
	// Output:
	// <label for="id_remember_me">Remember me:</label>
}

func ExampleField_CustomizeError() {
	fld := aform.Must(aform.NewCharField("Name", "", "", 10, 0))
	fld.CustomizeError(aform.ErrorWrapWithCode(fmt.Errorf("Please, enter enough characters"), aform.MinLengthErrorCode))
	fld.Clean("too_short")
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_name">Name</label>
	// <ul class="errorlist"><li id="err_0_id_name">Please, enter enough characters</li></ul>
	// <input type="text" name="name" value="too_short" id="id_name" minlength="10" aria-describedby="err_0_id_name" aria-invalid="true" required></div>
}

func ExampleField_AsDiv() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_name">Name</label><input type="text" name="name" id="id_name" maxlength="256" required></div>
}

func ExampleField_AsDiv_withHelpTextAndError() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fld.SetHelpText("Please, enter your name")
	fld.Clean("")
	fmt.Println(fld.AsDiv())
	// Output:
	// <div><label for="id_name">Name</label>
	// <ul class="errorlist"><li id="err_0_id_name">This field is required</li></ul>
	// <input type="text" name="name" id="id_name" maxlength="256" aria-describedby="helptext_id_name err_0_id_name" aria-invalid="true" required>
	// <span class="helptext" id="helptext_id_name">Please, enter your name</span></div>
}

func ExampleField_Name() {
	fld := aform.Must(aform.DefaultBooleanField("Remember me"))
	fmt.Println(fld.Name())
	// Output:
	// Remember me
}

func ExampleField_HTMLName() {
	fld := aform.Must(aform.DefaultBooleanField("Remember me"))
	fmt.Println(fld.HTMLName())
	// Output:
	// remember_me
}

func ExampleField_LabelTag() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fmt.Println(fld.LabelTag())
	// Output:
	// <label for="id_name">Name</label>
}

func ExampleField_LegendTag() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fmt.Println(fld.LegendTag())
	// Output:
	// <legend for="id_name">Name</legend>
}

func ExampleField_Widget() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fmt.Println(fld.Widget())
	// Output:
	// <input type="text" name="name" id="id_name" maxlength="256" required>
}

func ExampleField_Errors() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fld.Clean("")
	fmt.Println(fld.Errors())
	// Output:
	// <ul class="errorlist"><li id="err_0_id_name">This field is required</li></ul>
}

func ExampleField_Errors_withoutErrors() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fld.Clean("Paul")
	fmt.Printf("It renders the empty string: \"%s\"", fld.Errors())
	// Output:
	// It renders the empty string: ""
}

func ExampleField_HelpText() {
	fld := aform.Must(aform.DefaultCharField("Name"))
	fld.SetHelpText("Please, enter your name")
	fmt.Println(fld.HelpText())
	// Output:
	// <span class="helptext" id="helptext_id_name">Please, enter your name</span>
}
