package selection

// DefaultTemplate defines the default appearance of the selection and can
// be copied as a starting point for a custom template.

const DefaultTemplate = `
{{- if .Prompt -}}
  {{ .Prompt | faint | cyan }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- "⇣ " -}}
  {{- else -}}
    {{- "  " -}}
  {{- end -}}

  {{- if eq $.SelectedIndex $i }}
   {{- print (Foreground "32" (Bold "▸ ")) (Selected $choice) "\n" }}
  {{- else }}
    {{- print "  " (Unselected $choice) "\n" }}
  {{- end }}
{{- end}}`

// DefaultResultTemplate defines the default appearance with which the
// finale result of the selection is presented.
const DefaultResultTemplate = `
	{{- print .Prompt " " (Final .FinalChoice) "\n" -}}
`