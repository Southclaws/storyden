import { docs, blog as blogPosts } from "@/.source";
import { loader } from "fumadocs-core/source";
import { createMDXSource } from "fumadocs-mdx";
import { createOpenAPI, attachFile } from "fumadocs-openapi/server";

export const source = loader({
  baseUrl: "/docs",
  source: docs.toFumadocsSource(),
  pageTree: {
    attachFile,
  },
});

export const blog = loader({
  baseUrl: "/blog",
  source: createMDXSource(blogPosts),
});

export const openapi = createOpenAPI({
  disablePlayground: true,
});
