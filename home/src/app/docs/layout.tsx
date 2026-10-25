import { DocsLayout, DocsLayoutProps } from "fumadocs-ui/layouts/docs";
import type { ReactNode } from "react";

import { baseOptions, githubLink } from "@/app/layout.config";
import { source } from "@/lib/source";

import "fumadocs-openapi/css/preset.css";

const docsOptions: DocsLayoutProps = {
  ...baseOptions,
  tree: source.pageTree,
  links: [
    { text: "Storyden home", url: "/" },
    { text: "Live demo", url: "https://makeroom.club", external: true },
    githubLink,
  ],
};

export default async function Layout({ children }: { children: ReactNode }) {
  return <DocsLayout {...docsOptions}>{children}</DocsLayout>;
}
