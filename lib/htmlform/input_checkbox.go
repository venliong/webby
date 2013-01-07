package htmlform

import (
	"bytes"
)

type InputCheckbox struct {
	Name      string
	Value     string
	Id        string
	Class     string
	Selected  bool
	Mandatory bool
	error     error
	lang      Lang
}

func (fo *InputCheckbox) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="checkbox" name="{{.Name}}" value="{{.Value}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}{{if .Selected}}checked{{end}} />`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *InputCheckbox) Validate(values Values, files FileHeaders, single bool) error {
	if !values.Exist(fo.Name) {
		if fo.Mandatory {
			return FormError(fo.lang["ErrMandatoryCheckbox"])
		}
		return nil
	}

	fo.Selected = false

	for _, v := range values[fo.Name] {
		if v == fo.Value {
			fo.Selected = true
			break
		}
	}

	if !fo.Mandatory {
		return nil
	}

	if !fo.Selected {
		return FormError(fo.lang["ErrMandatoryCheckbox"])
	}

	return nil
}

func (fo *InputCheckbox) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputCheckbox) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputCheckbox) GetName() string {
	return fo.Name
}

func (fo *InputCheckbox) SetError(err error) {
	fo.error = err
}

func (fo *InputCheckbox) GetError() error {
	return fo.error
}

func (fo *InputCheckbox) GetStruct() FormHandler {
	return fo
}

func (fo *InputCheckbox) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *InputCheckbox) GetLang() Lang {
	return fo.lang
}