import * as Sentry from "@sentry/hono/node";
import type { NodeOptions } from "@sentry/node";

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  skipOpenTelemetrySetup: true,
  tracesSampleRate: undefined,
  dataCollection: {},
} as NodeOptions); // sentry/hono/node には skipOpenTelemetrySetup が存在しない
