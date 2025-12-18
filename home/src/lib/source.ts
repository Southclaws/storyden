import { blog as blogPosts, docs } from "@/.source/server";
import { loader } from "fumadocs-core/source";
import { toFumadocsSource } from "fumadocs-mdx/runtime/server";
import { openapiPlugin } from "fumadocs-openapi/server";
import { openapi } from "./openapi";

export const source = loader({
  baseUrl: "/docs",
  source: docs.toFumadocsSource(),
  plugins: [
    openapiPlugin(),
    {
      name: "openapi-toc",
      transformPageTree: {
        file(node, filePath) {
          if (!filePath) return node;

          const file = this.storage.read(filePath);
          if (!file || file.format !== "page") return node;

          const data = file.data as any;
          if (data._openapi?.toc) {
            data.toc = data._openapi.toc;
          }

          return node;
        },
      },
    },
  ],
});

export const blog = loader({
  baseUrl: "/blog",
  source: toFumadocsSource(blogPosts, []),
});

export { openapi };
