package htmlform

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type InputFile struct {
	Name      string
	Id        string
	Class     string
	Mime      []string
	Mandatory bool
	// In Byte
	Size    int64
	SizeErr string
	error   error
	lang    Lang
}

func init() {
	gob.Register(InputFile{})
}

func (fo *InputFile) Render(buf *bytes.Buffer) {
	const htmlstr = `<input type="file" name="{{.Name}}" {{if .IsId}}id="{{.Id}}" {{end}}{{if .IsClass}}class="{{.Class}}" {{end}}{{if .IsMime}}accept="{{.Accepted}}" {{end}}/>`
	if fo.error != nil {
		htmlRender(buf, fo.lang["ErrorTemplate"], fo.error.Error())
	}
	htmlRender(buf, htmlstr, fo)
}

func (fo *InputFile) Validate(values Values, files FileHeaders, single bool) error {
	if !files.Exist(fo.Name) {
		if fo.Mandatory {
			return FormError(fo.lang["ErrFileRequired"])
		}
	}

	if single {
		goto mime_check
	}

	if fo.Size <= 0 {
		goto mime_check
	}

	if files.GetSize(fo.Name) > fo.Size {
		return FormError(fo.SizeErr)
	}

mime_check:

	if len(fo.Mime) <= 0 {
		if fo.Mandatory {
			return FormError(fo.lang["ErrFileRequired"])
		}
		return nil
	}

	for _, value := range fo.Mime {
		if value == files.GetContentType(fo.Name) {
			return nil
		}
	}

	mimes := ""
	if len(fo.Mime) == 1 {
		mimes = fo.Mime[0]
	} else {
		mimes = strings.Join(fo.Mime[:len(fo.Mime)-1], ", ") + " & " + fo.Mime[len(fo.Mime)-1]
	}

	return FormError(fo.lang["ErrMimeCheck"] + mimes)
}

func (fo *InputFile) GetName() string {
	return fo.Name
}

func (fo *InputFile) SetError(err error) {
	fo.error = err
}

func (fo *InputFile) GetError() error {
	return fo.error
}

func (fo *InputFile) IsId() bool {
	return len(fo.Id) > 0
}

func (fo *InputFile) IsMime() bool {
	return len(fo.Mime) > 0
}

func (fo *InputFile) Accepted() string {
	return strings.Join(fo.Mime, "|")
}

func (fo *InputFile) IsClass() bool {
	return len(fo.Class) > 0
}

func (fo *InputFile) GetStruct() FormHandler {
	return fo
}

func (fo *InputFile) SetLang(lang Lang) {
	fo.lang = lang
}

func (fo *InputFile) GetLang() Lang {
	return fo.lang
}

func (fo InputFile) Eval() FormHandler {
	return &fo
}
