import { defineConfig } from "orval";

export default defineConfig({
  storyden: {
    input: {
      target: "../api/openapi.yaml",

      validation: true,
    },
    output: {
      target: "./src/api/openapi",
      mode: "tags",
      client: "swr",
      clean: true,
      prettier: true,
      schemas: "src/api/openapi/schemas",
      override: {
        mutator: {
          path: "./src/api/client.ts",
          name: "fetcher",
        },
      },
    },
  },
});
