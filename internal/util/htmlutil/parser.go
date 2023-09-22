package htmlutil

import (
	"bytes"
	"html/template"
	"path"
)

func ParseTemplate(templatePath string, data interface{}) (msg string, err error) {
	templ := template.New(path.Base(templatePath))
	templ, err = templ.ParseFiles(templatePath)
	if err != nil {
		return
	}

	tempBuffer := &bytes.Buffer{}

	err = templ.Execute(tempBuffer, data)
	if err != nil {
		return
	}

	return tempBuffer.String(), nil
}
