// Package greeting предназначен для вывода информации о сборке
package greeting

import (
	// init embed package
	_ "embed"
	"io"
	"text/template"
)

//go:embed greeting.txt
var greeting string

type buildInfo struct {
	Version string
	Date    string
	Commit  string
}

// PrintBuildInfo печатает информацию о сборке
func PrintBuildInfo(w io.Writer, version, date, commit string) error {
	info := buildInfo{
		Version: "N/A",
		Date:    "N/A",
		Commit:  "N/A",
	}

	if version != "" {
		info.Version = version
	}
	if date != "" {
		info.Date = date
	}
	if commit != "" {
		info.Commit = commit
	}

	tmpl := template.Must(template.New("greeting").Parse(greeting))
	return tmpl.Execute(w, info)
}
