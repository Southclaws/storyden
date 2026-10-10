import { canonicalURL, getIndexPage } from "@/lib/sitemap/public-index";

export const dynamic = "force-dynamic";

const TRAINING_AGENTS = [
  "GPTBot",
  "ClaudeBot",
  "CCBot",
  "Bytespider",
  "Google-Extended",
  "Applebot-Extended",
];

export async function GET(): Promise<Response> {
  try {
    const [threads, library] = await Promise.all([
      getIndexPage("threads", 1),
      getIndexPage("library", 1),
    ]);
    const searchEnabled = process.env["STORYDEN_SEARCH_CRAWLING"] !== "false";
    const trainingEnabled =
      process.env["STORYDEN_TRAINING_CRAWLING"] === "true";
    const lines = ["User-agent: *"];
    if (!searchEnabled || (threads === null && library === null)) {
      lines.push("Disallow: /");
    } else {
      lines.push(
        "Allow: /",
        "Disallow: /admin/",
        "Disallow: /settings/",
        "Disallow: /login",
        "Disallow: /register",
        "Disallow: /invitation/",
        "Disallow: /robots/",
        "Disallow: /trails/",
      );
      if (threads === null) lines.push("Disallow: /t/", "Disallow: /d/");
      if (library === null) lines.push("Disallow: /l/");
      lines.push(`Sitemap: ${canonicalURL("/sitemap.xml")}`);
    }
    if (!trainingEnabled) {
      for (const agent of TRAINING_AGENTS)
        lines.push("", `User-agent: ${agent}`, "Disallow: /");
    }
    return new Response(`${lines.join("\n")}\n`, {
      headers: {
        "Content-Type": "text/plain; charset=utf-8",
        "Cache-Control": "no-store",
      },
    });
  } catch (error) {
    console.error("Failed to determine public crawler access", error);
    return new Response("Crawler policy temporarily unavailable", {
      status: 503,
      headers: { "Cache-Control": "no-store" },
    });
  }
}
