import { InputOptions, OutputOptions, defineConfig } from "orval";

const input = {
  target: "../api/openapi.yaml",
  validation: true,
} as InputOptions;

const common = {
  mode: "tags" as const,
  clean: true,
  formatter: "prettier",
} as OutputOptions;

export default defineConfig({
  client: {
    input,
    output: {
      ...common,
      target: "src/api/openapi-client",
      client: "swr",
      schemas: "src/api/openapi-schema",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
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
        fetch: {
          includeHttpResponseReturnType: true,
          forceSuccessResponse: true,
        },
        mutator: {
          path: "./src/api/server.ts",
          name: "fetcher",
        },
      },
    },
  },
});
