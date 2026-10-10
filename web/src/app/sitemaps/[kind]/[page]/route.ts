import {
  XML_HEADERS,
  getIndexPage,
  urlsetXML,
} from "@/lib/sitemap/public-index";

export const dynamic = "force-dynamic";

export async function GET(
  _request: Request,
  context: { params: Promise<{ kind: string; page: string }> },
): Promise<Response> {
  const { kind, page: pageName } = await context.params;
  if (kind !== "threads" && kind !== "library")
    return new Response("Not found", { status: 404 });
  if (!/^[1-9]\d*\.xml$/.test(pageName))
    return new Response("Not found", { status: 404 });
  const page = Number(pageName.slice(0, -4));
  if (!Number.isSafeInteger(page) || page > 1_000_000)
    return new Response("Not found", { status: 404 });

  try {
    const result = await getIndexPage(kind, page);
    if (result === null || page > result.totalPages) {
      return new Response("Not found", {
        status: 404,
        headers: { "Cache-Control": "no-store" },
      });
    }
    return new Response(urlsetXML(kind, result.entries), {
      headers: XML_HEADERS,
    });
  } catch (error) {
    console.error("Failed to build sitemap page", error);
    return new Response("Sitemap temporarily unavailable", {
      status: 503,
      headers: { "Cache-Control": "no-store" },
    });
  }
}
