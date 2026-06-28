import * as Sentry from "@sentry/hono/node";
Sentry.init({
  dsn: process.env.SENTRY_DSN,
  tracesSampleRate: 0.0,
  dataCollection: {},
});
