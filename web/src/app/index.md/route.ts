import type {
  CategoryListResult,
  Info,
  NodeListResult,
  NodeWithAncestors,
  ThreadListResult,
} from "@/api/openapi-schema";
import {
  errorResponse,
  htmlToPublicMarkdown,
  inline,
  markdownResponse,
  pageNumber,
  publicJSON,
  publicURL,
} from "@/lib/markdown/public";
import { DefaultFeedConfig, FeedConfigSchema } from "@/lib/settings/feed";

export const dynamic = "force-dynamic";

export async function GET(request: Request) {
  try {
    const page = pageNumber(request);
    const info = await publicJSON<Info>("/info");
    const parsed = FeedConfigSchema.safeParse(info.metadata?.["feed"]);
    const feed = parsed.success ? parsed.data : DefaultFeedConfig;
    const blocks = feed.blocks;
    const canonical = publicURL(info.web_address, "/");
    const canonicalPage =
      page === 1 ? canonical : publicURL(info.web_address, `/?page=${page}`);
    const markdownPage = publicURL(
      info.web_address,
      `/index.md${page === 1 ? "" : `?page=${page}`}`,
    );
    const lines = [
      `# ${inline(info.title)}`,
      "",
      `Canonical: ${canonicalPage}`,
      "",
      inline(info.description),
    ];

    if (blocks.some((block) => block.type === "content")) {
      lines.push(
        "",
        htmlToPublicMarkdown(info.content, {
          web: info.web_address,
          api: info.api_address,
          base: canonical,
        }),
      );
    }

    if (blocks.some((block) => block.type === "categories")) {
      const categories = await publicJSON<CategoryListResult>("/categories");
      lines.push("", "## Categories", "");
      for (const category of categories.categories) {
        lines.push(
          `- [${inline(category.name)}](${publicURL(info.web_address, `/d/${encodeURIComponent(category.slug)}`)}): ${inline(category.description)}`,
        );
      }
    }

    const threadsBlock = blocks.find((block) => block.type === "threads");
    if (!threadsBlock && page > 1) return markdownResponse("# 404\n", 404);
    if (threadsBlock?.type === "threads") {
      const params = new URLSearchParams({ page: String(page) });
      if (threadsBlock.source === "uncategorised")
        params.set("categories", "null");
      const threads = await publicJSON<ThreadListResult>(`/threads?${params}`);
      if (page > 1 && page > threads.total_pages)
        return markdownResponse("# 404\n", 404);
      lines.push("", "## Discussions", "");
      for (const thread of threads.threads.filter(
        (item) => item.visibility === "published",
      )) {
        const link = publicURL(
          info.web_address,
          `/t/${encodeURIComponent(thread.slug)}.md`,
        );
        lines.push(
          `- [${inline(thread.title)}](${link}) — ${inline(thread.author.name)}, ${thread.createdAt}`,
        );
      }
      lines.push("", `Page ${page} of ${Math.max(threads.total_pages, 1)}`);
      if (page > 1)
        lines.push(
          `Previous: ${publicURL(info.web_address, `/index.md?page=${page - 1}`)}`,
        );
      if (threads.next_page) {
        lines.push(
          `Next: ${publicURL(info.web_address, `/index.md?page=${threads.next_page}`)}`,
        );
      }
    }

    const libraryBlock = blocks.find((block) => block.type === "library");
    if (libraryBlock?.type === "library") {
      let nodes: NodeListResult["nodes"];
      let nextLibraryPage: string | undefined;
      if (libraryBlock.node) {
        const parent = await publicJSON<NodeWithAncestors>(
          `/nodes/${encodeURIComponent(libraryBlock.node)}`,
        );
        if (parent.visibility === "published") {
          const children = await publicJSON<NodeListResult>(
            `/nodes/${encodeURIComponent(parent.slug)}/children?page=1`,
          );
          nodes = children.nodes;
          if (children.next_page) {
            nextLibraryPage = publicURL(
              info.web_address,
              `/l/${encodeURIComponent(parent.slug)}.md?page=${children.next_page}`,
            );
          }
        } else {
          nodes = [];
        }
      } else {
        nodes = (await publicJSON<NodeListResult>("/nodes?depth=0")).nodes;
      }
      lines.push("", "## Library", "");
      for (const node of nodes.filter(
        (item) => item.visibility === "published",
      )) {
        lines.push(
          `- [${inline(node.name)}](${publicURL(info.web_address, `/l/${encodeURIComponent(node.slug)}.md`)}): ${inline(node.description)}`,
        );
      }
      if (nextLibraryPage) lines.push(`More Library pages: ${nextLibraryPage}`);
    }

    return markdownResponse(lines.join("\n"), 200, canonicalPage, markdownPage);
  } catch (error) {
    return errorResponse(error);
  }
}
