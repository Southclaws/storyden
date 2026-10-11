import type { Info, Thread } from "@/api/openapi-schema";
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

type Context = { params: Promise<{ slug: string }> };

export async function GET(request: Request, context: Context) {
  try {
    const page = pageNumber(request);
    const { slug } = await context.params;
    if (!slug || slug.length > 200) return markdownResponse("# 404\n", 404);

    const [info, thread] = await Promise.all([
      publicJSON<Info>("/info"),
      publicJSON<Thread>(`/threads/${encodeURIComponent(slug)}?page=${page}`),
    ]);
    if (thread.visibility !== "published")
      return markdownResponse("# 404\n", 404);
    if (page > 1 && page > thread.replies.total_pages) {
      return markdownResponse("# 404\n", 404);
    }

    const canonical = publicURL(
      info.web_address,
      `/t/${encodeURIComponent(thread.slug)}`,
    );
    const canonicalPage = page === 1 ? canonical : `${canonical}?page=${page}`;
    const markdownPage = `${canonical}.md${page === 1 ? "" : `?page=${page}`}`;
    const lines = [
      `# ${inline(thread.title)}`,
      "",
      `Canonical: ${canonicalPage}`,
      `Author: ${inline(thread.author.name)} (@${inline(thread.author.handle)})`,
      `Published: ${thread.createdAt}`,
      `Updated: ${thread.updatedAt}`,
      "",
      htmlToPublicMarkdown(thread.body, {
        web: info.web_address,
        api: info.api_address,
        base: canonical,
      }),
    ];

    if (thread.replies.replies.length) {
      lines.push("", "## Replies", "");
      for (const reply of thread.replies.replies) {
        if (reply.visibility !== "published") continue;
        lines.push(
          `### ${inline(reply.author.name)} — ${reply.createdAt}`,
          "",
          htmlToPublicMarkdown(reply.body, {
            web: info.web_address,
            api: info.api_address,
            base: canonical,
          }),
          "",
        );
      }
    }

    lines.push(`Page ${page} of ${Math.max(thread.replies.total_pages, 1)}`);
    if (page > 1) lines.push(`Previous: ${canonical}.md?page=${page - 1}`);
    if (thread.replies.next_page) {
      lines.push(`Next: ${canonical}.md?page=${thread.replies.next_page}`);
    }
    return markdownResponse(lines.join("\n"), 200, canonicalPage, markdownPage);
  } catch (error) {
    return errorResponse(error);
  }
}
