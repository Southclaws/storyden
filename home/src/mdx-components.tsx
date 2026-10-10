import defaultComponents from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";

import { Mermaid } from "@/components/Mermaid";

export function getMDXComponents(): MDXComponents {
  return {
    ...defaultComponents,
    Mermaid,
  };
}
