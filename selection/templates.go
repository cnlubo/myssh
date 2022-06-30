package selection

// DefaultTemplate defines the default appearance of the selection and can
// be copied as a starting point for a custom template.
//{{ .HeadPrompt | bold | cyan }}
const DefaultTemplate = `
{{- if .Prompt -}}
  {{ .Prompt | faint | cyan  }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}

{{- if eq $.SelectedIndex $i }}
  {{- print | reset (green  (Cursor) " ") (Selected $choice) "\n" }}
{{- else }}
   {{- print | reset " " ( faint (Unselected $choice)) "\n" }}
{{- end }}

{{- end}}`

//const DefaultTemplate = `
//{{- if .Prompt -}}
//  {{ .Prompt | faint | cyan  }}
//{{ end -}}
//{{ if .IsFiltered }}
//  {{- print .FilterPrompt " " .FilterInput }}
//{{ end }}
//
//{{- range  $i, $choice := .Choices }}
//  {{- if IsScrollUpHintPosition $i }}
//    {{- "⇡ " -}}
//  {{- else if IsScrollDownHintPosition $i -}}
//    {{- "⇣ " -}}
//  {{- else -}}
//    {{- "  " -}}
//  {{- end -}}
//
//{{- if eq $.SelectedIndex $i }}
//  {{- print | reset (green  "» ") (Selected $choice) "\n" }}
//{{- else }}
//   {{- print | reset " " (Unselected $choice) "\n" }}
//{{- end }}

//{{- end}}`

// DefaultResultTemplate defines the default appearance with which the
// finale result of the selection is presented.
const DefaultResultTemplate = `
	{{- print .Prompt " " (Final .FinalChoice) "\n" -}}
`