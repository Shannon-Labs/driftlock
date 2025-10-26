{{/*
Expand the name of the chart.
*/}}
{{- define "driftlock.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "driftlock.fullname" -}}
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
{{- define "driftlock.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "driftlock.labels" -}}
helm.sh/chart: {{ include "driftlock.chart" . }}
{{ include "driftlock.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "driftlock.selectorLabels" -}}
app.kubernetes.io/name: {{ include "driftlock.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "driftlock.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "driftlock.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Kafka brokers connection string
*/}}
{{- define "driftlock.kafka.brokers" -}}
{{- if .Values.kafka.enabled }}
{{- $brokers := list -}}
{{- range $i := until (int .Values.kafka.replicaCount) -}}
{{- $brokers = append $brokers (printf "%s-kafka-%d.%s-kafka-headless.%s.svc.cluster.local:9092" (include "driftlock.fullname" .) $i (include "driftlock.fullname" .) .Release.Namespace) -}}
{{- end -}}
{{- join "," $brokers -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end }}

{{/*
Redis connection URL
*/}}
{{- define "driftlock.redis.url" -}}
{{- if .Values.redis.enabled }}
{{- if .Values.redis.auth.enabled }}
{{- printf "redis://:%s@%s-redis-master.%s.svc.cluster.local:6379" .Values.redis.auth.password (include "driftlock.fullname" .) .Release.Namespace }}
{{- else }}
{{- printf "redis://%s-redis-master.%s.svc.cluster.local:6379" (include "driftlock.fullname" .) .Release.Namespace }}
{{- end }}
{{- else }}
{{- "" }}
{{- end }}
{{- end }}

{{/*
ClickHouse connection URL
*/}}
{{- define "driftlock.clickhouse.url" -}}
{{- if .Values.clickhouse.enabled }}
{{- if .Values.clickhouse.auth.enabled }}
{{- printf "clickhouse://default:%s@%s-clickhouse.%s.svc.cluster.local:9000/%s" .Values.clickhouse.auth.password (include "driftlock.fullname" .) .Release.Namespace .Values.clickhouse.configmap.db }}
{{- else }}
{{- printf "clickhouse://%s-clickhouse.%s.svc.cluster.local:9000/%s" (include "driftlock.fullname" .) .Release.Namespace .Values.clickhouse.configmap.db }}
{{- end }}
{{- else }}
{{- "" }}
{{- end }}
{{- end }}
