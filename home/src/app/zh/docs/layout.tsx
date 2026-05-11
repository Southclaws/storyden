import { baseOptionsZh } from "@/app/layout.config";
import { sourceZh } from "@/lib/source";
import { DocsLayout, DocsLayoutProps } from "fumadocs-ui/layouts/docs";
import type { ReactNode } from "react";

const docsOptions: DocsLayoutProps = {
  ...baseOptionsZh,
  tree: sourceZh.pageTree,
  links: [],
};

export default async function Layout({ children }: { children: ReactNode }) {
  return <DocsLayout {...docsOptions}>{children}</DocsLayout>;
}
