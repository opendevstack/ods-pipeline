package rendermanifest

import (
	"io"
	"text/template"
)

// RenderManifest renders the given template with the passed data,
// writing the result to w.
func RenderManifest(w io.Writer, tmpl *template.Template, data map[string]string) error {
	if _, err := w.Write(
		[]byte("# File is generated; DO NOT EDIT.\n\n"),
	); err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
