import { Type, type TypeBoxTypeProvider } from "@fastify/type-provider-typebox";
// Import the framework and instantiate it
import Fastify, {
  errorCodes,
  type FastifyError,
  type FastifyInstance,
  type FastifyPluginOptions,
  type FastifyReply,
  type FastifyRequest,
  type FastifySchema,
} from "fastify";
import pretty from "pino-pretty";
import {
  CustomFastifyError,
  type FastifyReplyTypebox,
  type FastifyRequestTypebox,
  type FastifyTypebox,
  NotFoundError,
} from "./types";

function plainRoutes(fastify: FastifyInstance) {
  fastify.get("/plan/1", async function handler(_request, _reply) {
    return { hello: "world" };
  });

  fastify.route({
    method: "GET",
    url: "/plan/2",
    handler: async function handler(_request, _reply) {
      return { hello: "world" };
    },
  });
}

function typeboxRoutes(fastify: FastifyTypebox) {
  const schema = {
    body: Type.Object({
      x: Type.String(),
      y: Type.Number(),
      z: Type.Boolean(),
    }),
  } satisfies FastifySchema;

  fastify.post("/typebox", { schema: schema }, (req) => {
    const { x, y, z } = req.body;
    return { message: `x: ${x}, y: ${y}, z: ${z}` };
  });
}

function typeboxWithHandlerRoutes(fastify: FastifyTypebox) {
  const schema = {
    params: Type.Object({
      id: Type.Number(),
    }),
  } satisfies FastifySchema;

  function handler(
    req: FastifyRequestTypebox<typeof schema>,
    _reply: FastifyReplyTypebox<typeof schema>,
  ) {
    const { id } = req.params;
    return { message: `id: ${id}` };
  }

  fastify.get("/typebox/:id", { schema }, handler);
}

function withPrefixRoutes(
  fastify: FastifyTypebox,
  _opts: FastifyPluginOptions,
  done: (err?: Error) => void,
) {
  fastify.get("/", (_req) => {
    return "ok";
  });
  done();
}

function customErrorRoutes(fastify: FastifyTypebox) {
  fastify.get("/error/custom", (_req, _reply) => {
    throw new CustomFastifyError("something");
  });
  fastify.get("/error/notfound", (_req, _reply) => {
    throw new NotFoundError("nowhere", { resource: "nobody" });
  });
  fastify.get("/error/string", (_req, _reply) => {
    throw "error!";
  });
  fastify.get("/error/error", (_req, _reply) => {
    throw new Error("aaa");
  });
}

type ErrorBody = {
  code?: string;
  error?: string;
  message: string;
  statusCode: number;
  [key: string]: unknown;
};

function makeErrorResult(error: unknown): ErrorBody | Error {
  if (error instanceof NotFoundError) {
    return {
      error: "Resource Not Found",
      message: error.message,
      statusCode: 404,
      resource: error.resource,
    };
  }
  if (error instanceof Error) {
    return error;
  }
  return {
    message: String(error),
    statusCode: 500,
  };
}

async function errorHandler(
  error: FastifyError | Error,
  _req: FastifyRequest,
  reply: FastifyReply,
) {
  const errorResult = makeErrorResult(error);
  if (errorResult instanceof Error) {
    await reply.send(errorResult);
  } else {
    await reply.code(errorResult.statusCode).send(errorResult);
  }
}

async function main() {
  const fastify = Fastify({
    logger: {
      stream: pretty(),
    },
  }).withTypeProvider<TypeBoxTypeProvider>();
  fastify.setErrorHandler(errorHandler);

  fastify.errorHandler;

  plainRoutes(fastify);
  typeboxRoutes(fastify);
  typeboxWithHandlerRoutes(fastify);
  customErrorRoutes(fastify);
  fastify.register(withPrefixRoutes, { prefix: "/prefix" });

  // Run the fastify!
  try {
    await fastify.listen({ port: 3000 });
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
}

await main();
