import "./instrument.ts";

import { serve } from "@hono/node-server";
import { httpInstrumentationMiddleware } from "@hono/otel";
import { sentry } from "@sentry/hono/node";
import { Hono } from "hono";
import { logger } from "hono/logger";
import * as z from "zod";

const NEXT_URL = process.env.NEXT_URL ?? "http://localhost:3000";

const ResultSchema = z.object({
  result: z.number(),
});

const app = new Hono();
app.use(sentry(app));
app.use(httpInstrumentationMiddleware());
app.use(logger());

app.get("/sum/:n", async (c) => {
  const n = parseInt(c.req.param("n"), 10);
  if (isNaN(n) || n < 0) {
    return c.json({ error: "n must be a non-negative integer" }, 422);
  }
  if (n === 0) {
    return c.json({ result: 0 });
  }

  const resp = await fetch(`${NEXT_URL}/sum/${n - 1}`);
  if (!resp.ok) {
    return c.json({ error: `upstream error: ${resp.status}` }, 502);
  }
  const data = ResultSchema.parse(await resp.json());
  return c.json({ result: data.result + n });
});

app.get("/error", async (_c) => {
  throw new Error("something went wrong");
});

serve(
  {
    fetch: app.fetch,
    port: 3000,
  },
  (info) => {
    console.log(`Server is running on http://localhost:${info.port}`);
  },
);
