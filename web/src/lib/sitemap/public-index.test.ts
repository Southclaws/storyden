import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  escapeXML,
  getIndexPage,
  sitemapIndexXML,
  urlsetXML,
} from "./public-index";

vi.mock("@/config", () => ({
  WEB_ADDRESS: "https://community.example",
  getAPIAddress: () => "https://api.community.example",
}));

beforeEach(() => vi.restoreAllMocks());

describe("anonymous sitemap source", () => {
  it("fetches published threads without forwarding an admin session", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(
        JSON.stringify({
          threads: [
            {
              slug: "one",
              updatedAt: "2026-01-02T03:04:05Z",
              visibility: "published",
            },
          ],
          total_pages: 2,
        }),
        { status: 200 },
      ),
    );

    const page = await getIndexPage("threads", 1);
    expect(page?.entries).toEqual([
      { slug: "one", updatedAt: "2026-01-02T03:04:05Z" },
    ]);
    expect(fetchMock).toHaveBeenCalledOnce();
    const [url, options] = fetchMock.mock.calls[0]!;
    expect(String(url)).toBe(
      "https://api.community.example/api/threads?page=1&visibility=published&ignore_pinned=true",
    );
    expect(options).toMatchObject({
      headers: { Accept: "application/json" },
      cache: "no-store",
    });
    expect(options?.signal).toBeInstanceOf(AbortSignal);
  });

  it("uses the newest visible thread activity as lastmod", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(
        JSON.stringify({
          threads: [
            {
              slug: "one",
              updatedAt: "2026-01-01T00:00:00Z",
              last_reply_at: "2026-01-03T00:00:00Z",
              visibility: "published",
            },
          ],
          total_pages: 1,
        }),
        { status: 200 },
      ),
    );
    expect((await getIndexPage("threads", 1))?.entries[0]?.updatedAt).toBe(
      "2026-01-03T00:00:00Z",
    );
  });

  it("never indexes a non-published thread even if the API returns one", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(
        JSON.stringify({
          threads: [
            {
              slug: "draft",
              updatedAt: "2026-01-02T03:04:05Z",
              visibility: "draft",
            },
          ],
          total_pages: 1,
        }),
        { status: 200 },
      ),
    );
    await expect(getIndexPage("threads", 1)).rejects.toThrow("non-published");
  });

  it("treats guest denial separately from upstream failure", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch");
    fetchMock.mockResolvedValueOnce(new Response(null, { status: 401 }));
    expect(await getIndexPage("library", 1)).toBeNull();
    fetchMock.mockResolvedValueOnce(new Response("error", { status: 500 }));
    await expect(getIndexPage("library", 1)).rejects.toThrow("500");
  });
});

describe("sitemap XML", () => {
  it("escapes canonical URLs, includes lastmod, and lists every page", () => {
    expect(escapeXML(`&<>"'`)).toBe("&amp;&lt;&gt;&quot;&apos;");
    const index = sitemapIndexXML({ threads: 2, library: 1 });
    expect(index).toContain("https://community.example/sitemaps/threads/2.xml");
    expect(index.match(/<sitemap>/g)).toHaveLength(3);

    const xml = urlsetXML("library", [
      { slug: `a&b"<`, updatedAt: "2026-01-02T03:04:05Z" },
    ]);
    expect(xml).toContain("https://community.example/l/a%26b%22%3C");
    expect(xml).toContain("<lastmod>2026-01-02T03:04:05.000Z</lastmod>");
  });

  it("rejects an index with more than 50,000 sitemap entries", () => {
    expect(() => sitemapIndexXML({ threads: 50_001, library: 0 })).toThrow(
      "limit",
    );
  });
});
