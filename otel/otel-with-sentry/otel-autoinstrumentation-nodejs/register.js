"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

const opentelemetry = require("@opentelemetry/sdk-node");
const { diag, DiagConsoleLogger } = require("@opentelemetry/api");
const { getStringFromEnv, diagLogLevelFromString } = require("@opentelemetry/core");
const { getNodeAutoInstrumentations, getResourceDetectorsFromEnv } = require("./utils");

const { getMetricReader } = require("./custom");

const logLevel = getStringFromEnv("OTEL_LOG_LEVEL");
if (logLevel != null) {
  diag.setLogger(new DiagConsoleLogger(), {
    logLevel: diagLogLevelFromString(logLevel),
  });
}

const sdk = new opentelemetry.NodeSDK({
  instrumentations: getNodeAutoInstrumentations(),
  resourceDetectors: getResourceDetectorsFromEnv(),

  // NodeSDK does not read OTEL_EXPORTER_OTLP_METRICS_DEFAULT_HISTOGRAM_AGGREGATION (issue #3920).
  // So, we detect the aggregation preference from the environment and pass it to the exporter explicitly.
  metricReader: getMetricReader(),
});

try {
  sdk.start();
  diag.info("OpenTelemetry automatic instrumentation started successfully");
} catch (error) {
  diag.error("Error initializing OpenTelemetry SDK", error);
}

async function shutdown() {
  try {
    await sdk.shutdown();
  } catch (error) {
    diag.error("Error terminating OpenTelemetry SDK", error);
  }
}

// Gracefully shutdown SDK if a SIGTERM is received
process.on("SIGTERM", shutdown);
// Gracefully shutdown SDK if Node.js is exiting normally
process.once("beforeExit", shutdown);
