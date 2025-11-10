{{/*
Define the conveyor.namespace template if set with namespace or .Release.Namespace is set
*/}}
{{- define "conveyor.namespace" -}}
  {{- default .Release.Namespace .Values.namespace -}}
{{- end }}