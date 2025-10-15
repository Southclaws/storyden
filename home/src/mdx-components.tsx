import { openapi } from "@/lib/openapi";
import { APIPage } from "fumadocs-openapi/ui";
import defaultComponents from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";
import { Mermaid } from "@/components/Mermaid";

export function getMDXComponents(): MDXComponents {
  return {
    ...defaultComponents,
    APIPage: (props) => <APIPage {...openapi.getAPIPageProps(props)} />,
    Mermaid,
  };
}
