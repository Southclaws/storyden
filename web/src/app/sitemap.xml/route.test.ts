import { DefaultDiscoverySettings } from "@/lib/settings/discovery";
import { getDiscoverySettings } from "@/lib/settings/discovery-server";

vi.mock("@/lib/settings/discovery-server", () => ({
  getDiscoverySettings: vi.fn(),
}));

beforeEach(() => {
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
  });
});

import { beforeEach, expect, it, vi } from "vitest";

import { getIndexPage } from "@/lib/sitemap/public-index";

import { GET } from "./route";

vi.mock("@/lib/sitemap/public-index", () => ({
  getIndexPage: vi.fn(),
  sitemapIndexXML: ({
    threads,
    library,
  }: {
    threads: number;
    library: number;
  }) => `${threads}:${library}`,
  XML_HEADERS: { "Content-Type": "application/xml" },
}));

beforeEach(() => {
  vi.mocked(getIndexPage).mockResolvedValue({ entries: [], totalPages: 2 });
});

it("lists only public content types available to guests", async () => {
  vi.mocked(getIndexPage).mockResolvedValueOnce(null);
  const response = await GET();
  expect(response.status).toBe(200);
  expect(await response.text()).toBe("0:2");
});

it("returns 404 when neither public content type is available", async () => {
  vi.mocked(getIndexPage).mockResolvedValue(null);
  expect((await GET()).status).toBe(404);
});

it("returns 503 on backend failure", async () => {
  vi.mocked(getIndexPage).mockRejectedValue(new Error("unavailable"));
  vi.spyOn(console, "error").mockImplementation(() => {});
  expect((await GET()).status).toBe(503);
});

it("does not serve sitemaps when search indexing is disabled", async () => {
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
    search_indexing: false,
  });
  const response = await GET();
  expect(response.status).toBe(404);
  expect(response.headers.get("Cache-Control")).toBe("no-store");
  expect(getIndexPage).not.toHaveBeenCalled();
});
