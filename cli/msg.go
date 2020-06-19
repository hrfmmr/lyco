package cli

import (
	"fmt"
	"io"
	"text/template"

	"github.com/hrfmmr/lyco"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var TemplateData map[string]interface{}

var Templates = template.Must(template.New("Help").Parse(`
{{define "NameVer"}}lyco {{.Version}}{{end}}

{{define "OneLine"}}A terminal user interface for pomodoro technique🍅{{end}}

{{define "Header"}}{{template "NameVer" .}}

{{template "OneLine"}}
{{end}}
`))

func init() {
	TemplateData = map[string]interface{}{
		"Version": lyco.Version,
	}
}

func WriteHelp(p *flags.Parser, w io.Writer) {
	if err := Templates.ExecuteTemplate(w, "Header", TemplateData); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w)
	p.WriteHelp(w)
}
