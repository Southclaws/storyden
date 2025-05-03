import defaultComponents from "fumadocs-ui/mdx";
import { APIPage } from "fumadocs-openapi/ui";
import { openapi } from "@/lib/source";
import type { MDXComponents } from "mdx/types";

export function getMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultComponents,
    APIPage: (props) => {
      return (
        <div>
          <APIPage {...openapi.getAPIPageProps(props)} />
        </div>
      );
    },
    ...components,
  };
}
