package confirmation

// TemplateArrow is a template where the current choice is indicated by an
// arrow.

const TemplateArrow = `
{{- Bold .Prompt -}}
{{ if .YesSelected -}}
	{{- print (Bold " ▸Yes ") " No" -}}
{{- else if .NoSelected -}}
	{{- print "  Yes " (Bold "▸No") -}}
{{- else -}}
	{{- "  Yes  No" -}}
{{- end -}}
`

// ResultTemplateArrow is the ResultTemplate that matches TemplateArrow.
const ResultTemplateArrow = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- Foreground "32" "Yes" -}}
{{- else -}}
	{{- Foreground "32" "No" -}}
{{- end }}
`

const TemplateYN = `
{{- .Prompt | bold | red -}}
{{ if .YesSelected  -}}
	{{- print | red (bold " [Y/") red ("n") red (bold "]") -}}
{{- else if .NoSelected -}}
	{{- print | red (bold " [y/N]") -}}
{{- else -}}
	{{- " [y/n]" -}}
{{- end -}}
`
const ResultTemplateYN = `
{{- .Prompt | bold | green -}}
{{ if .FinalValue -}}
    {{- print | green (bold " [Y/n]") -}}
{{- else -}}
	{{- print | green (bold " [y/N]") -}}
{{- end }}
`