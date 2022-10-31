package t_executor

import (
	"errors"
	"fmt"
	"io"
	"html/template"
)

func ExecuteTemplate(tmpls *template.Template, tmplName string, w io.Writer, data interface{}) error {
	var err error
	layout := tmpls.Lookup("layout")
	if layout == nil {
		return errors.New("no layout template")
	}

	layout, err = layout.Clone()
	if err != nil {
		return err
	}

	t := tmpls.Lookup(tmplName)
	if t == nil {
		return fmt.Errorf("template not found: %s", tmplName)
	}

	_, err = layout.AddParseTree("", t.Tree)
	if err != nil {
		return err
	}

	return layout.Execute(w, data)
}
