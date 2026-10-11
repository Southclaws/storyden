"use server";

import { canonicalURL, getIndexPage } from "@/lib/sitemap/public-index";

const TRAINING_AGENTS = [
  "GPTBot",
  "ClaudeBot",
  "CCBot",
  "Bytespider",
  "Google-Extended",
  "Applebot-Extended",
];

export async function generateRobotsTxt(
  searchIndexing: boolean,
  trainingCrawling: boolean,
): Promise<string> {
  const lines = ["User-agent: *"];

  if (searchIndexing) {
    const [threads, library] = await Promise.all([
      getIndexPage("threads", 1),
      getIndexPage("library", 1),
    ]);

    lines.push(
      "Allow: /",
      "Disallow: /admin/",
      "Disallow: /settings/",
      "Disallow: /login$",
      "Disallow: /register$",
      "Disallow: /invitation/",
      "Disallow: /robots/",
      "Disallow: /trails/",
    );

    if (threads === null) {
      lines.push("Disallow: /t/", "Disallow: /d/");
    }

    if (library === null) {
      lines.push("Disallow: /l/");
    }

    if (threads !== null || library !== null) {
      lines.push(`Sitemap: ${canonicalURL("/sitemap.xml")}`);
    }
  } else {
    lines.push("Disallow: /");
  }

  if (!trainingCrawling) {
    for (const agent of TRAINING_AGENTS) {
      lines.push("", `User-agent: ${agent}`, "Disallow: /");
    }
  }

  return `${lines.join("\n")}\n`;
}
