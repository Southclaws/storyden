package robot

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed global.md.tmpl
var globalInstructionTemplateSource string

//go:embed identity.md.tmpl
var identityInstructionTemplateSource string

var globalInstructionTemplate = template.Must(
	template.New("global-instruction").
		Option("missingkey=error").
		Parse(globalInstructionTemplateSource),
)

var identityInstructionTemplate = template.Must(
	template.New("identity-instruction").
		Option("missingkey=error").
		Parse(identityInstructionTemplateSource),
)

func renderInstructionTemplate(tmpl *template.Template, data any) (string, error) {
	return renderNamedInstructionTemplate(tmpl, "", data)
}

func renderNamedInstructionTemplate(tmpl *template.Template, name string, data any) (string, error) {
	var output bytes.Buffer
	var err error
	if name == "" {
		err = tmpl.Execute(&output, data)
	} else {
		err = tmpl.ExecuteTemplate(&output, name, data)
	}
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}
