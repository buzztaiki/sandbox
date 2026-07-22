"use strict";

const { diag } = require("@opentelemetry/api");
const { getStringFromEnv } = require("@opentelemetry/core");


const { OTLPMetricExporter: OTLPGrpcMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-grpc');
const { OTLPMetricExporter: OTLPHttpMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-http');
const { OTLPMetricExporter: OTLPProtoMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-proto');
const { PeriodicExportingMetricReader, AggregationType, InstrumentType } = require("@opentelemetry/sdk-metrics");

function getOtlpMetricExporterFromEnv() {
  const protocol = (
    getStringFromEnv("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL") ?? getStringFromEnv("OTEL_EXPORTER_OTLP_PROTOCOL")
  )?.trim() || "http/protobuf";
  const aggregationPreference = getMetricAggregationPreferenceFromEnv();

  switch (protocol) {
  case "grpc":
    return new OTLPGrpcMetricExporter({ aggregationPreference });
  case "http/json":
    return new OTLPHttpMetricExporter({ aggregationPreference });
  case "http/protobuf":
    return new OTLPProtoMetricExporter({ aggregationPreference });
  default:
    diag.warn(`Unsupported OTLP metrics protocol: "${protocol}". Using http/protobuf.`);
    return new OTLPProtoMetricExporter({ aggregationPreference });
  }
}

function getMetricAggregationPreferenceFromEnv() {
  const aggregation = (
    getStringFromEnv("OTEL_EXPORTER_OTLP_METRICS_DEFAULT_HISTOGRAM_AGGREGATION")
  )?.trim() || "explicit_bucket_histogram";
  const aggregationPreference = (type) => (instrumentType) =>
        instrumentType === InstrumentType.HISTOGRAM ? { type } : { type: AggregationType.DEFAULT };

  switch (aggregation) {
  case "base2_exponential_bucket_histogram":
    return aggregationPreference(AggregationType.EXPONENTIAL_HISTOGRAM);
  case "explicit_bucket_histogram":
    return aggregationPreference(AggregationType.EXPLICIT_BUCKET_HISTOGRAM);
  default:
    diag.warn(`Unsupported histogram aggregation: "${aggregation}". Using explicit_bucket_histogram.`);
    return aggregationPreference(AggregationType.EXPLICIT_BUCKET_HISTOGRAM);
  }
}

function getMetricReader() {
  return new PeriodicExportingMetricReader({ exporter: getOtlpMetricExporterFromEnv() });
}

module.exports = { getMetricReader };
