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

{{/*
Returns the LoadBalancer hostname from the ingress controller service.
*/}}
{{- define "baseDomain" -}}
{{- $svc := (lookup "v1" "Service" "kube-system" "traefik") -}}
{{- if $svc -}}
  {{- $annotations := $svc.metadata.annotations -}}
  {{- if $annotations -}}
    {{- $loadbalancerId := index $annotations "kubernetes.civo.com/loadbalancer-id" -}}
    {{- if $loadbalancerId -}}
      {{- printf "%s.lb.civo.com" $loadbalancerId -}}
    {{- else -}}
      {{- .Values.baseDomain -}}
    {{- end -}}
  {{- else -}}
    {{- .Values.baseDomain -}}
  {{- end -}}
{{- else -}}
  {{- .Values.baseDomain -}}
{{- end -}}
{{- end -}}