import type {
  Info,
  NodeListResult,
  NodeWithAncestors,
} from "@/api/openapi-schema";
import { parseNodeMetadata } from "@/lib/library/metadata";
import {
  errorResponse,
  htmlToPublicMarkdown,
  inline,
  markdownResponse,
  pageNumber,
  publicJSON,
  publicURL,
} from "@/lib/markdown/public";

export const dynamic = "force-dynamic";

type Context = { params: Promise<{ slug: string[] }> };

export async function GET(request: Request, context: Context) {
  try {
    const page = pageNumber(request);
    const { slug } = await context.params;
    if (
      !slug?.length ||
      slug.length > 12 ||
      slug.some((part) => !part || part.length > 200)
    ) {
      return markdownResponse("# 404\n", 404);
    }
    const target = slug.at(-1)!;
    const [info, node] = await Promise.all([
      publicJSON<Info>("/info"),
      publicJSON<NodeWithAncestors>(`/nodes/${encodeURIComponent(target)}`),
    ]);
    if (node.visibility !== "published")
      return markdownResponse("# 404\n", 404);

    const canonical = publicURL(
      info.web_address,
      `/l/${encodeURIComponent(node.slug)}`,
    );
    const canonicalPage = page === 1 ? canonical : `${canonical}?page=${page}`;
    const markdownPage = `${canonical}.md${page === 1 ? "" : `?page=${page}`}`;
    const lines = [
      `# ${inline(node.name)}`,
      "",
      `Canonical: ${canonicalPage}`,
      `Author: ${inline(node.owner.name)} (@${inline(node.owner.handle)})`,
      `Published: ${node.createdAt}`,
      `Updated: ${node.updatedAt}`,
      "",
      inline(node.description),
      "",
      htmlToPublicMarkdown(node.content ?? "", {
        web: info.web_address,
        api: info.api_address,
        base: canonical,
      }),
    ];

    const hasDirectory = parseNodeMetadata(node.meta).layout?.blocks.some(
      (block) => block.type === "directory",
    );
    const shouldListChildren =
      page > 1 ||
      node.children.length > 0 ||
      (node.hide_child_tree && hasDirectory);
    let childPage: NodeListResult | undefined;
    if (shouldListChildren) {
      childPage = await publicJSON<NodeListResult>(
        `/nodes/${encodeURIComponent(node.slug)}/children?page=${page}`,
      );
      if (page > 1 && page > childPage.total_pages)
        return markdownResponse("# 404\n", 404);
    }
    const children =
      childPage?.nodes.filter((child) => child.visibility === "published") ??
      [];
    if (children.length || childPage) {
      lines.push("", "## Child pages", "");
      for (const child of children) {
        lines.push(
          `- [${inline(child.name)}](${publicURL(info.web_address, `/l/${encodeURIComponent(child.slug)}.md`)})`,
        );
      }
      lines.push("", `Page ${page} of ${Math.max(childPage!.total_pages, 1)}`);
      if (page > 1) lines.push(`Previous: ${canonical}.md?page=${page - 1}`);
      if (childPage?.next_page)
        lines.push(`Next: ${canonical}.md?page=${childPage.next_page}`);
    }
    return markdownResponse(lines.join("\n"), 200, canonicalPage, markdownPage);
  } catch (error) {
    return errorResponse(error);
  }
}
