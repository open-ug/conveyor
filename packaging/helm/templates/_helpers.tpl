{{- define "conveyor.name" -}}
{{ include "conveyor.fullname" . }}
{{- end }}

{{- define "conveyor.fullname" -}}
{{- if .Values.nameOverride }}
{{ .Values.nameOverride }}
{{- else }}
{{ .Chart.Name }}
{{- end }}
{{- end }}

{{- define "conveyor.namespace" -}}
{{- if .Values.namespace }}
{{- .Values.namespace }}
{{- else }}
{{- .Release.Namespace }}
{{- end }}
{{- end }}
