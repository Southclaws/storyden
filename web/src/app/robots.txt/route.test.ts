import { afterEach, beforeEach, expect, it, vi } from "vitest";

import { getIndexPage } from "@/lib/sitemap/public-index";

import { GET } from "./route";

vi.mock("@/lib/sitemap/public-index", () => ({
  getIndexPage: vi.fn(),
  canonicalURL: (path: string) => `https://community.example${path}`,
}));

beforeEach(() => {
  vi.mocked(getIndexPage).mockResolvedValue({ entries: [], totalPages: 1 });
});
afterEach(() => vi.unstubAllEnvs());

it("allows search by default while disallowing known training crawlers", async () => {
  const response = await GET();
  const text = await response.text();
  expect(response.status).toBe(200);
  expect(text).toContain("Sitemap: https://community.example/sitemap.xml");
  expect(text).toContain("User-agent: GPTBot\nDisallow: /");
});

it("honours separate operator crawl preferences", async () => {
  vi.stubEnv("STORYDEN_SEARCH_CRAWLING", "false");
  vi.stubEnv("STORYDEN_TRAINING_CRAWLING", "true");
  const text = await (await GET()).text();
  expect(text).toContain("User-agent: *\nDisallow: /");
  expect(text).not.toContain("User-agent: GPTBot");
  expect(text).not.toContain("Sitemap:");
});

it("disallows crawling if guests cannot list either public content type", async () => {
  vi.mocked(getIndexPage).mockResolvedValue(null);
  expect(await (await GET()).text()).toContain("User-agent: *\nDisallow: /");
});

it("returns an error when the API is unavailable", async () => {
  vi.mocked(getIndexPage).mockRejectedValue(new Error("unavailable"));
  vi.spyOn(console, "error").mockImplementation(() => {});
  expect((await GET()).status).toBe(503);
});
