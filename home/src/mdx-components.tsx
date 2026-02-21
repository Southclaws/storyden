import { openapi } from "@/lib/openapi";
import { createAPIPage } from "fumadocs-openapi/ui";
import defaultComponents from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";
import { Mermaid } from "@/components/Mermaid";
import type { MediaAdapter } from "fumadocs-openapi";

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

const APIPage = createAPIPage(openapi, {
  playground: { enabled: false },
  mediaAdapters,
});

export function getMDXComponents(): MDXComponents {
  return {
    ...defaultComponents,
    APIPage,
    Mermaid,
  };
}
