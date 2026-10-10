import {
  XML_HEADERS,
  getIndexPage,
  sitemapIndexXML,
} from "@/lib/sitemap/public-index";

export const dynamic = "force-dynamic";

export async function GET(): Promise<Response> {
  try {
    const [threads, library] = await Promise.all([
      getIndexPage("threads", 1),
      getIndexPage("library", 1),
    ]);
    if (threads === null && library === null) {
      return new Response("Public indexing is unavailable", {
        status: 404,
        headers: { "Cache-Control": "no-store" },
      });
    }
    const xml = sitemapIndexXML({
      threads: threads?.totalPages ?? 0,
      library: library?.totalPages ?? 0,
    });
    return new Response(xml, { headers: XML_HEADERS });
  } catch (error) {
    console.error("Failed to build sitemap index", error);
    return new Response("Sitemap temporarily unavailable", {
      status: 503,
      headers: { "Cache-Control": "no-store" },
    });
  }
}
