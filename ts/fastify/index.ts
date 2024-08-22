// Import the framework and instantiate it
import Fastify, { type FastifyInstance } from "fastify";

function helloRoutes(fastify: FastifyInstance) {
  // Declare a route
  fastify.get("/", async function handler(_request, _reply) {
    return { hello: "world" };
  });
}

async function main() {
  const fastify = Fastify({
    logger: true,
  });

  helloRoutes(fastify);

  // Run the server!
  try {
    await fastify.listen({ port: 3000 });
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
}

await main();
