package htmlform

import (
	"fmt"
	"github.com/CJ-Jackson/webby"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

func TestTextarea(t *testing.T) {
	fmt.Println("Textarea Test:\r\n")

	form := New(nil,
		&Textarea{
			Name:    "text",
			MinChar: 1,
			MaxChar: 8,
		},
	)

	fmt.Println(form.Render())
	fmt.Println()

	web := &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"text": []string{"hello"},
			},
			MultipartForm: &multipart.Form{},
		},
	}

	if !form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()

	web = &webby.Web{
		Req: &http.Request{
			Form: url.Values{
				"text": []string{"hellohello"},
			},
			MultipartForm: &multipart.Form{},
		},
	}

	if form.IsValid(web) {
		t.Fail()
	}

	fmt.Println(form.Render())
	fmt.Println()
}
