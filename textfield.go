package fork

import (
	"fmt"
	"net/http"
	"net/mail"
	"strings"
)

func textwidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="text" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func newtextfield(name string, widget Widget, validaters []interface{}, filters []interface{}) Field {
	return &textfield{
		name: name,
		Text: "",
		processor: NewProcessor(widget,
			validaters,
			filters,
		),
	}
}

func TextField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return newtextfield(name, textwidget(options...), validaters, filters)
}

type textfield struct {
	name         string
	Text         string
	validateable bool
	*processor
}

func (t *textfield) New() Field {
	var newfield textfield = *t
	t.validateable = false
	return &newfield
}

func (t *textfield) Name(name ...string) string {
	if len(name) > 0 {
		t.name = strings.Join(name, "-")
	}
	return t.name
}

func (t *textfield) Get() *Value {
	return NewValue(t.Text)
}

func (t *textfield) Set(r *http.Request) {
	v := t.Filter(t.Name(), r)
	t.Text = v.String()
	t.validateable = true
}

func (t *textfield) Validateable() bool {
	return t.validateable
}

func textareawidget(options ...string) Widget {
	return NewWidget(WithOptions(`<textarea name="{{ .Name }}" %s>{{ .Text }}</textarea>`, options...))
}

func TextAreaField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return newtextfield(name, textareawidget(options...), validaters, filters)
}

func hiddenwidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="hidden" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func HiddenField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return newtextfield(name, hiddenwidget(options...), validaters, filters)
}

func passwordwidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="password" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func PassWordField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return newtextfield(name, passwordwidget(options...), validaters, filters)
}

func emailwidget(options ...string) Widget {
	return NewWidget(WithOptions(`<input type="email" name="{{ .Name }}" value="{{ .Text }}" %s>`, options...))
}

func EmailField(name string, validaters []interface{}, filters []interface{}, options ...string) Field {
	return newtextfield(name, emailwidget(options...), append(validaters, ValidEmail), nil)
}

func ValidEmail(t *textfield) error {
	if t.validateable {
		_, err := mail.ParseAddress(t.Text)
		if err != nil {
			return fmt.Errorf("Invalid email address: %s", err.Error())
		}
	}
	return nil
}
