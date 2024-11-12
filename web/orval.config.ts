import { defineConfig } from "orval";

const input = {
  target: "../api/openapi.yaml",
  validation: false,
};

const common = {
  mode: "tags" as const,
  clean: true,
  prettier: true,
};

export default defineConfig({
  client: {
    input,
    output: {
      ...common,
      target: "src/api/openapi-client",
      client: "swr",
      schemas: "src/api/openapi-schema",
      override: {
        mutator: {
          path: "./src/api/client.ts",
          name: "fetcher",
        },
      },
    },
  },
  server: {
    input,
    output: {
      ...common,
      target: "src/api/openapi-server",
      schemas: "src/api/openapi-schema",
      client: "fetch",
      override: {
        mutator: {
          path: "./src/api/server.ts",
          name: "fetcher",
        },
      },
    },
  },
});
