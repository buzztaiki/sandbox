import { Type, type TypeBoxTypeProvider } from "@fastify/type-provider-typebox";
// Import the framework and instantiate it
import Fastify, {
  type FastifyInstance,
  type FastifyPluginOptions,
  type FastifySchema,
} from "fastify";
import type {
  FastifyReplyTypebox,
  FastifyRequestTypebox,
  FastifyTypebox,
} from "./types";

function plainRoutes(fastify: FastifyInstance) {
  // Declare a route
  fastify.get("/", async function handler(_request, _reply) {
    return { hello: "world" };
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

function withPrefix(
  fastify: FastifyTypebox,
  _opts: FastifyPluginOptions,
  done: (err?: Error) => void,
) {
  fastify.get("/", (_req) => {
    return "ok";
  });
  done();
}

async function main() {
  const fastify = Fastify({
    logger: true,
  }).withTypeProvider<TypeBoxTypeProvider>();

  plainRoutes(fastify);
  typeboxRoutes(fastify);
  typeboxWithHandlerRoutes(fastify);
  fastify.register(withPrefix, { prefix: "/prefix" });

  // Run the fastify!
  try {
    await fastify.listen({ port: 3000 });
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
}

await main();
