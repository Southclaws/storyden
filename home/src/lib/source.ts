import { blog as blogPosts, docs } from "@/.source";
import { loader } from "fumadocs-core/source";
import { createMDXSource } from "fumadocs-mdx";
import { attachFile, openapiPlugin } from "fumadocs-openapi/server";
import { openapi } from "./openapi";

export const source = loader({
  baseUrl: "/docs",
  source: docs.toFumadocsSource(),
  pageTree: {
    attachFile,
  },
  plugins: [
    openapiPlugin(),
    {
      // hack to fix openapi pages not having table of contents.
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
  source: createMDXSource(blogPosts),
});

export { openapi };
