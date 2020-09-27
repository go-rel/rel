{{with .PDoc}}
{{if $.IsMain}}
> {{ base .ImportPath }}
{{comment_md .Doc}}
{{else}}

# {{ .Name }}

`import "{{.ImportPath}}"`

## Overview

{{comment_md .Doc}}

{{example_html $ ""}}

{{with .Consts}}## Constants
{{range .}}{{node $ .Decl | pre}}
{{comment_md .Doc}}{{end}}{{end}}

{{with .Vars}}## Variables
{{range .}}{{node $ .Decl | pre}}
{{comment_md .Doc}}{{end}}{{end}}
{{range .Funcs}}{{$name_html := html .Name}}## func {{$name_html}}
{{node $ .Decl | pre}}
{{comment_md .Doc}}
{{example_html $ .Name}}
{{callgraph_html $ "" .Name}}{{end}}


{{range .Types}}{{$tname := .Name}}{{$tname_html := html .Name}}## type {{$tname_html}}
{{node $ .Decl | pre}}
{{comment_md .Doc}}{{range .Consts}}
{{node $ .Decl | pre }}
{{comment_md .Doc}}{{end}}{{range .Vars}}
{{node $ .Decl | pre }}
{{comment_md .Doc}}{{end}}
{{example_html $ $tname}}
{{implements_html $ $tname}}
{{methodset_html $ $tname}}

{{range .Funcs}}{{$name_html := html .Name}}### func {{$name_html}}
{{node $ .Decl | pre}}
{{comment_md .Doc}}
{{example_html $ .Name}}{{end}}
{{callgraph_html $ "" .Name}}

{{range .Methods}}{{$name_html := html .Name}}### func ({{md .Recv}}) {{$name_html}}
{{node $ .Decl | pre}}
{{comment_md .Doc}}
{{$name := printf "%s_%s" $tname .Name}}{{example_html $ $name}}
{{callgraph_html $ .Recv .Name}}
{{end}}{{end}}{{end}}
{{with $.Notes}}
{{range $marker, $content := .}}

{{end}}
{{end}}
{{end}}
