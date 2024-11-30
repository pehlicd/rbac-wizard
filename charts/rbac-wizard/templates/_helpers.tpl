{{/*
Expand the name of the chart.
*/}}
{{- define "rbac-wizard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rbac-wizard.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "rbac-wizard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "rbac-wizard.labels" -}}
helm.sh/chart: {{ include "rbac-wizard.chart" . }}
{{ include "rbac-wizard.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "rbac-wizard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "rbac-wizard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "rbac-wizard.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "rbac-wizard.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "rbac-wizard.clusterRoleName" -}}
{{- if .Values.clusterRole.create }}
{{- default (include "rbac-wizard.fullname" .) .Values.clusterRole.name }}
{{- else }}
{{- default "default" .Values.clusterRole.name }}
{{- end }}
{{- end }}

{{- define "rbac-wizard.clusterRoleBindingName" -}}
{{- if .Values.clusterRoleBinding.create }}
{{- default (include "rbac-wizard.fullname" .) .Values.clusterRoleBinding.name }}
{{- else }}
{{- default "default" .Values.clusterRoleBinding.name }}
{{- end }}
{{- end }}
