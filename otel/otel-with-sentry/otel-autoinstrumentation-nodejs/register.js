"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

const opentelemetry = require("@opentelemetry/sdk-node");
const api = require("@opentelemetry/api");
const { getStringFromEnv, diagLogLevelFromString } = require("@opentelemetry/core");
const { getNodeAutoInstrumentations, getResourceDetectorsFromEnv } = require("./utils");

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
            api.diag.warn(`Unsupported OTLP metrics protocol: "${protocol}". Using http/protobuf.`);
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
            api.diag.warn(`Unsupported histogram aggregation: "${aggregation}". Using explicit_bucket_histogram.`);
            return aggregationPreference(AggregationType.EXPLICIT_BUCKET_HISTOGRAM);
    }
}


const logLevel = getStringFromEnv("OTEL_LOG_LEVEL");
if (logLevel != null) {
    api.diag.setLogger(new api.DiagConsoleLogger(), {
        logLevel: diagLogLevelFromString(logLevel),
    });
}

const sdk = new opentelemetry.NodeSDK({
    instrumentations: getNodeAutoInstrumentations(),
    resourceDetectors: getResourceDetectorsFromEnv(),

    // NodeSDK does not read OTEL_EXPORTER_OTLP_METRICS_DEFAULT_HISTOGRAM_AGGREGATION (issue #3920).
    // So, we detect the aggregation preference from the environment and pass it to the exporter explicitly.
    metricReader: new PeriodicExportingMetricReader({ exporter: getOtlpMetricExporterFromEnv() }),
});

try {
    sdk.start();
    api.diag.info("OpenTelemetry automatic instrumentation started successfully");
} catch (error) {
    api.diag.error("Error initializing OpenTelemetry SDK", error);
}

async function shutdown() {
    try {
        await sdk.shutdown();
    } catch (error) {
        api.diag.error("Error terminating OpenTelemetry SDK", error);
    }
}

// Gracefully shutdown SDK if a SIGTERM is received
process.on("SIGTERM", shutdown);
// Gracefully shutdown SDK if Node.js is exiting normally
process.once("beforeExit", shutdown);
