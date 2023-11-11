package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/BelyaevEI/shortener/internal/app"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// BuildData хранит в себе информацию о текущей версии, дате и значении коммита
type BuildData struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

// Template - переменная, содержит в себе темплейт для stdout с информацией о текущей сборке
const Template = `	Build version: {{if .BuildVersion}} {{.BuildVersion}} {{else}} N/A {{end}}
	Build version: {{if .BuildDate}} {{.BuildDate}} {{else}} N/A {{end}}
	Build version: {{if .BuildCommit}} {{.BuildCommit}} {{else}} N/A {{end}}
`

func main() {
	buildInfo()
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}

func buildInfo() {
	d := &BuildData{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	t := template.Must(template.New("buildTags").Parse(Template))
	err := t.Execute(os.Stdout, *d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}
}
