import "./instrument.ts";

import * as Sentry from "@sentry/node";
import express, { type NextFunction, type Request, type Response } from "express";
import { pinoHttp } from "pino-http";
import z from "zod";

const NEXT_URL = process.env.NEXT_URL ?? "http://localhost:3001";

const ResultSchema = z.object({
  result: z.number(),
});

const app = express();
app.use(pinoHttp());

app.get("/sum/:n", async (req, res) => {
  const n = parseInt(req.params.n, 10);
  if (isNaN(n) || n < 0) {
    res.status(422).json({ error: "n must be a non-negative integer" });
    return;
  }
  if (n === 0) {
    n
    res.json({ result: 0 });
    return;
  }

  const resp = await fetch(`${NEXT_URL}/sum/${n - 1}`);
  if (!resp.ok) {
    res.status(502).json({ error: `upstream error: ${resp.status}` });
    return;
  }
  const data = ResultSchema.parse(await resp.json());
  res.json({ result: data.result + n });
});

app.get("/error", (_req, _res) => {
  throw new Error("something went wrong");
});

Sentry.setupExpressErrorHandler(app);
app.use((err: unknown, _req: Request, res: Response, next: NextFunction) => {
  if (res.headersSent) {
    return next(err);
  }
  res.status(500).json({ error: "Internal Server Error" });
});

app.listen(3001);
