import defaultComponents from "fumadocs-ui/mdx";
import { openapi } from "@/lib/source";
import type { MDXComponents } from "mdx/types";

export function getMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultComponents,
    APIPage: openapi.APIPage,
    ...components,
  };
}
