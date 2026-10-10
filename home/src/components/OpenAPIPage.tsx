"use client";

import type { MediaAdapter } from "fumadocs-openapi";
import { createOpenAPIPage } from "fumadocs-openapi/ui";

const mediaAdapters: Record<string, MediaAdapter> = {
  "text/plain": {
    encode(data) {
      return String(data.body ?? "");
    },
    generateExample(data, ctx) {
      const bodyStr = JSON.stringify(data.body ?? "", null, 2);
      if (ctx.lang === "js") {
        return `const body = ${bodyStr};`;
      }
      return bodyStr;
    },
  },
  "application/zip": {
    encode(data) {
      return data.body as BodyInit;
    },
    generateExample() {
      // not supported
      return undefined;
    },
  },
};

export const OpenAPIPage = createOpenAPIPage({
  playground: { enabled: false },
  mediaAdapters,
});
