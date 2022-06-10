package confirmation

// TemplateArrow is a template where the current choice is indicated by an
// arrow.

const TemplateArrow = `
{{- .Prompt | bold | cyan -}}
{{ if .YesSelected -}}
	{{- print (bold (yellow " ▸Yes ")) (cyan " No") -}}
{{- else if .NoSelected -}}
	{{- print (cyan "  Yes ") (bold (yellow "▸No")) -}}
{{- else -}}
	{{- "  Yes  No" -}}
{{- end -}}
`
const ResultTemplateArrow = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- print | cyan (bold "Yes") -}}
{{- else -}}
	{{- print | cyan (bold "No") -}}
{{- end }}
`

const TemplateYN = `
{{- .Prompt | bold | cyan  -}}
{{ if .YesSelected  -}}
	{{- print (cyan " [") (bold (yellow "Y")) (cyan "/n]") -}}
{{- else if .NoSelected -}}
    {{- print (cyan " [y/") (bold (yellow "N")) (cyan "]") -}}
{{- else -}}
	{{- " [y/n]" -}}
{{- end -}}
`
const ResultTemplateYN = `
{{- .Prompt | bold | cyan -}}
{{ if .FinalValue -}}
    {{- print | cyan (bold " [Y/n]") -}}
{{- else -}}
	{{- print | cyan (bold " [y/N]") -}}
{{- end }}
`