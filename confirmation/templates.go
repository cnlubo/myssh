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

// TemplateYN is a classic template with ja [yn] indicator where the current
// value is capitalized and bold.
//{{- Bold .Prompt -}}
//const TemplateYN = `
//{{- .Prompt | Bold | cyan -}}
//{{ if .YesSelected  -}}
//	{{- print " [" (Bold "Y") "/n]" -}}
//{{- else if .NoSelected -}}
//	{{- print " red([y/" (Bold "N") "])" -}}
//{{- else -}}
//	{{- " [y/n]" -}}
//{{- end -}}
//`
const TemplateYN = `
{{- .Prompt | Bold | cyan -}}
{{ if .YesSelected  -}}
	{{- print "[" (Bold "Y") "/n]" -}}
{{- else if .NoSelected -}}
	{{- print (Bold " [") (Bold (cyan "y/")) (Bold (yellow "N")) (Bold "]") -}}
{{- else -}}
	{{- " [y/n] " | red -}}
{{- end -}}
`

// ResultTemplateYN is the ResultTemplate that matches TemplateYN.
const ResultTemplateYN = `
{{- .Prompt -}}
{{ if .FinalValue -}}
	{{- print " [" (Foreground "32" (Bold "Y")) "/n]" -}}
{{- else -}}
	{{- print " [y/" (Foreground "32" (Bold "N")) "]" -}}
{{- end }}
`