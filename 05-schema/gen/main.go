package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

func main() {
	pkg := flag.String("package", "", "package name")
	flag.Parse()

	schema := Schema(
		Message("Person",
			Field("Name", "string"),
			Field("Age", "uint64"),
			Field("Path", "Points"),
		),

		SliceOf("Points", "Point"),
		Message("Point",
			Field("X", "uint64"),
			Field("Y", "uint64"),
		),
	)

	err := T.Execute(os.Stdout, map[string]any{
		"Package": *pkg,
		"Schema":  schema,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var T = template.Must(template.New("").Funcs(template.FuncMap{
	"native": func(v string) string {
		switch v {
		case "string":
			return "String"
		case "uint64":
			return "Uint64"
		}
		return ""
	},
}).Parse(`package {{.Package}}

import (
	"errors"

	"github.com/golang-estonia/structs-to-bytes/est"
	"github.com/zeebo/errs"
)

{{ range $m := .Schema.Messages -}}
{{- if $m.SliceOf }}
type {{$m.Name}} []{{$m.SliceOf}}

func (m {{$m.Name}}) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteUint64(uint64(len(m))); err != nil {
		return err
	}
	for i := range m {
		if err = stream.WriteMessage(m[i].EncodeEst); err != nil {
			return err
		}
	}
	return nil
}

func (m *{{$m.Name}}) DecodeEst(stream *est.Stream) (err error) {
	n, err := stream.ReadUint64()
	if err != nil {
		return err
	}
	*m = make({{$m.Name}}, int(n))
	for i := range *m {
		if err = stream.ReadMessage((*m)[i].DecodeEst); err != nil {
			return err
		}
	}
	return nil
}
{{ else }}
type {{$m.Name}} struct {
	{{- range $f := $m.Fields }}
	{{ $f.Name }} {{ $f.Type }}
	{{- end }}
}

func (m {{$m.Name}}) EncodeEst(stream *est.Stream) (err error) {
	{{- range $f := $m.Fields -}}
	{{- with $native := (native $f.Type) }}
	if err = stream.Write{{$native}}(m.{{$f.Name}}); err != nil {
		return errs.Wrap(err)
	}
	{{ else }}
	if err = stream.WriteMessage(m.{{$f.Name}}.EncodeEst); err != nil {
		return errs.Wrap(err)
	}
	{{ $f.Name }}
	{{ end -}}
	{{- end }}
	return nil
}

func (m *{{$m.Name}}) DecodeEst(stream *est.Stream) (err error) {
	{{- range $f := $m.Fields -}}
	{{- with $native := (native $f.Type) }}
	if m.{{$f.Name}}, err = stream.Read{{$native}}(); err != nil {
		return errs.Wrap(err)
	}
	{{ else }}
	if err = stream.ReadMessage(m.{{$f.Name}}.DecodeEst); err != nil {
		return errs.Wrap(err)
	}
	{{ $f.Name }}
	{{ end -}}
	{{- end }}
	return nil
}
{{ end -}}
{{- end -}}
`))
