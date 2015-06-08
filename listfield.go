package fork

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type listField struct {
	base   Field
	Fields []Field
	*baseField
	*processor
}

func (l *listField) New() Field {
	var newfield listField = *l
	newfield.base = l.base.New()
	newfield.validateable = false
	return &newfield
}

func (l *listField) Get() *Value {
	return NewValue(l.Fields)
}

func (lf *listField) Set(r *http.Request) {
	lf.Fields = nil
	pl := 0
	for k, _ := range r.PostForm {
		if getregexp(lf.name).MatchString(k) {
			pl++
		}
	}
	if pl > 0 {
		for x := 0; x < pl; x++ {
			nf := lf.base.New()
			renameField(lf.name, x, nf)
			nf.Set(r)
			lf.Fields = append(lf.Fields, nf)
		}
	}
	lf.validateable = true
}

func listFieldsWidget(options ...string) Widget {
	return NewWidget(WithOptions(`<fieldset name="{{ .Name }}" %s><ul>{{ range $x := .Fields }}<li>{{ .Render $x }}</li>{{ end }}</ul></fieldset>`, options...))
}

func renameField(name string, number int, field Field) Field {
	field.Name(name, strconv.Itoa(number), field.Name())
	return field
}

func renameFields(name string, number int, field Field) []Field {
	var ret []Field
	for i := 0; i < number; i++ {
		fd := field.New()
		ret = append(ret, renameField(name, i, fd))
	}
	return ret
}

func getregexp(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^%s-", name))
}

func ListField(name string, startwith int, starting Field, options ...string) Field {
	return &listField{
		base:      starting,
		Fields:    renameFields(name, startwith, starting),
		baseField: newBaseField(name),
		processor: NewProcessor(
			listFieldsWidget(options...),
			nilValidater,
			nilFilterer,
		),
	}
}
