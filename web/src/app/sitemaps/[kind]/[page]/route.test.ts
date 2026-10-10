import { beforeEach, expect, it, vi } from "vitest";

import { getIndexPage } from "@/lib/sitemap/public-index";

import { GET } from "./route";

vi.mock("@/lib/sitemap/public-index", () => ({
  getIndexPage: vi.fn(),
  urlsetXML: (_kind: string, entries: unknown[]) => `entries:${entries.length}`,
  XML_HEADERS: { "Content-Type": "application/xml" },
}));

const request = new Request("https://community.example/sitemaps/threads/1.xml");
const context = (kind: string, page: string) => ({
  params: Promise.resolve({ kind, page }),
});

beforeEach(() =>
  vi
    .mocked(getIndexPage)
    .mockResolvedValue({
      entries: [{ slug: "one", updatedAt: "2026-01-02T03:04:05Z" }],
      totalPages: 2,
    }),
);

it("serves a valid numbered page", async () => {
  const response = await GET(request, context("threads", "1.xml"));
  expect(response.status).toBe(200);
  expect(await response.text()).toBe("entries:1");
  expect(getIndexPage).toHaveBeenCalledWith("threads", 1);
});

it.each(["0.xml", "01.xml", "1", "3.xml", "999999999999999999.xml"])(
  "rejects an invalid or missing page %s",
  async (page) => {
    expect((await GET(request, context("threads", page))).status).toBe(404);
  },
);

it("rejects unknown sitemap kinds", async () => {
  expect((await GET(request, context("private", "1.xml"))).status).toBe(404);
});
