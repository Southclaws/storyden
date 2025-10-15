import { createOpenAPI } from "fumadocs-openapi/server";

export const openapi = createOpenAPI({
  // the OpenAPI schema, you can also give it an external URL.
  input: ["../api/openapi.yaml"],
  disablePlayground: true,
  mediaAdapters: {
    // Handle text/plain content type (used by /beacon endpoint)
    "text/plain": {
      encode(data) {
        return JSON.stringify(data.body);
      },
      generateExample(data, ctx) {
        if (ctx.lang === "js") {
          return `const body = ${JSON.stringify(data.body, null, 2)};`;
        }
        return JSON.stringify(data.body, null, 2);
      },
    },
  },
});
