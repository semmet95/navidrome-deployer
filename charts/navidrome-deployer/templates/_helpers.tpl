{{/*
Common labels
*/}}
{{- define "commonLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
App labels
*/}}
{{- define "appLabels" -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}