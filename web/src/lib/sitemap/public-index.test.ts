import { JSDOM } from "jsdom";
import { cookies } from "next/headers";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { getIndexPage, sitemapIndexXML, urlsetXML } from "./public-index";

vi.mock("@/config", () => ({
  WEB_ADDRESS: "https://community.example",
  getAPIAddress: () => "https://api.community.example",
}));

vi.mock("next/headers", () => ({
  cookies: vi.fn(),
  headers: vi.fn(),
}));

beforeEach(() => vi.restoreAllMocks());

describe("anonymous sitemap source", () => {
  it("fetches published threads without reading the visitor's session", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(
      Response.json({
        threads: [
          {
            slug: "one",
            updatedAt: "2026-01-02T03:04:05Z",
            visibility: "published",
          },
        ],
        total_pages: 2,
      }),
    );

    expect(await getIndexPage("threads", 1)).toEqual({
      entries: [{ slug: "one", updatedAt: "2026-01-02T03:04:05Z" }],
      totalPages: 2,
    });
    expect(fetchMock).toHaveBeenCalledOnce();

    const request = fetchMock.mock.calls[0]![0] as Request;
    expect(request.url).toBe(
      "https://api.community.example/api/threads?page=1&visibility=published&ignore_pinned=true",
    );
    expect(request.cache).toBe("no-store");
    expect(request.headers.has("Cookie")).toBe(false);
    expect(request.headers.has("Authorization")).toBe(false);
    expect(cookies).not.toHaveBeenCalled();
  });

  it.each([
    [undefined, "2026-01-02T00:00:00Z"],
    ["2026-01-01T00:00:00Z", "2026-01-02T00:00:00Z"],
    ["2026-01-03T00:00:00Z", "2026-01-03T00:00:00Z"],
  ])(
    "uses the latest thread activity when the last reply is %s",
    async (lastReplyAt, expected) => {
      vi.spyOn(globalThis, "fetch").mockResolvedValue(
        Response.json({
          threads: [
            {
              slug: "one",
              updatedAt: "2026-01-02T00:00:00Z",
              last_reply_at: lastReplyAt,
              visibility: "published",
            },
          ],
          total_pages: 1,
        }),
      );

      expect((await getIndexPage("threads", 1))?.entries[0]?.updatedAt).toBe(
        expected,
      );
    },
  );

  it("requests library pages from the API and retains its pagination metadata", async () => {
    const nodes = Array.from({ length: 101 }, (_, i) => ({
      slug: `node-${i + 1}`,
      updatedAt: "2026-01-02T03:04:05Z",
      visibility: "published",
    }));
    const fetchMock = vi
      .spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(
        Response.json({ nodes: nodes.slice(0, 100), total_pages: 2 }),
      )
      .mockResolvedValueOnce(
        Response.json({ nodes: nodes.slice(100), total_pages: 2 }),
      );

    const firstPage = await getIndexPage("library", 1);
    const secondPage = await getIndexPage("library", 2);

    expect(firstPage?.totalPages).toBe(2);
    expect(firstPage?.entries).toHaveLength(100);
    expect(firstPage?.entries[0]?.slug).toBe("node-1");
    expect(firstPage?.entries[99]?.slug).toBe("node-100");
    expect(secondPage).toEqual({
      entries: [{ slug: "node-101", updatedAt: "2026-01-02T03:04:05Z" }],
      totalPages: 2,
    });

    for (const [i, [input]] of fetchMock.mock.calls.entries()) {
      expect((input as Request).url).toBe(
        `https://api.community.example/api/nodes?page=${i + 1}&visibility=published&format=flat`,
      );
    }
  });

  it("omits library sitemap pages when the public node list is empty", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      Response.json({ nodes: [], total_pages: 0 }),
    );

    expect(await getIndexPage("library", 1)).toEqual({
      entries: [],
      totalPages: 0,
    });
  });

  it.each(["threads", "library"] as const)(
    "rejects non-published %s content",
    async (kind) => {
      vi.spyOn(globalThis, "fetch").mockResolvedValue(
        Response.json({
          [kind === "threads" ? "threads" : "nodes"]: [
            {
              slug: "draft",
              updatedAt: "2026-01-02T03:04:05Z",
              visibility: "draft",
            },
          ],
          total_pages: 1,
        }),
      );

      await expect(getIndexPage(kind, 1)).rejects.toThrow("non-published");
    },
  );

  it.each([401, 403])(
    "treats HTTP %s as unavailable to guests",
    async (status) => {
      vi.spyOn(globalThis, "fetch").mockResolvedValue(
        new Response(null, { status }),
      );

      expect(await getIndexPage("library", 1)).toBeNull();
    },
  );

  it("propagates upstream failures", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response("error", { status: 500 }),
    );

    await expect(getIndexPage("library", 1)).rejects.toMatchObject({
      status: 500,
    });
  });
});

describe("sitemap XML", () => {
  it("lists every sitemap page in a valid XML document", () => {
    const xml = sitemapIndexXML({ threads: 2, library: 1 });
    const document: Document = new JSDOM(xml, {
      contentType: "application/xml",
    }).window.document;

    expect(document.documentElement.localName).toBe("sitemapindex");
    expect(document.documentElement.namespaceURI).toBe(
      "http://www.sitemaps.org/schemas/sitemap/0.9",
    );
    expect(
      Array.from(
        document.querySelectorAll("sitemap > loc"),
        (node) => node.textContent,
      ),
    ).toEqual([
      "https://community.example/sitemaps/threads/1.xml",
      "https://community.example/sitemaps/threads/2.xml",
      "https://community.example/sitemaps/library/1.xml",
    ]);
  });

  it("encodes slugs safely and serializes lastmod as an ISO timestamp", () => {
    const xml = urlsetXML("library", [
      { slug: `a&b"<`, updatedAt: "2026-01-02T04:04:05+01:00" },
    ]);
    const document: Document = new JSDOM(xml, {
      contentType: "application/xml",
    }).window.document;

    expect(document.documentElement.localName).toBe("urlset");
    expect(document.querySelector("url > loc")?.textContent).toBe(
      "https://community.example/l/a%26b%22%3C",
    );
    expect(document.querySelector("url > lastmod")?.textContent).toBe(
      "2026-01-02T03:04:05.000Z",
    );
    expect(document.querySelectorAll("url")).toHaveLength(1);
  });

  it.each([50_001, -1, 1.5, NaN])(
    "rejects an invalid sitemap page count of %s",
    (threads) => {
      expect(() => sitemapIndexXML({ threads, library: 0 })).toThrow("limit");
    },
  );
});
