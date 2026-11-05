{{- define "self-node-remediation.agentImage" -}}
{{- if .Values.agent.image }}
{{- .Values.agent.image }}
{{- else }}
{{- .Values.manager.image.repository }}{{ if not (contains "@" .Values.manager.image.repository) }}:{{ .Values.manager.image.tag | default .Chart.AppVersion }}{{ end }}
{{- end }}
{{- end }}
