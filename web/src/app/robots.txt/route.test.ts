import { beforeEach, expect, it, vi } from "vitest";

import { DefaultDiscoverySettings } from "@/lib/settings/discovery";
import { getDiscoverySettings } from "@/lib/settings/discovery-server";
import { generateRobotsTxt } from "@/lib/settings/robots-txt";
import { getIndexPage } from "@/lib/sitemap/public-index";

import { GET } from "./route";

vi.mock("@/lib/settings/discovery-server", () => ({
  getDiscoverySettings: vi.fn(),
}));

vi.mock("@/lib/sitemap/public-index", () => ({
  getIndexPage: vi.fn(),
  canonicalURL: (path: string) => `https://community.example${path}`,
}));

beforeEach(() => {
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
  });
  vi.mocked(getIndexPage).mockResolvedValue({ entries: [], totalPages: 1 });
});

it("allows search by default while disallowing known training crawlers", async () => {
  const response = await GET();
  const text = await response.text();
  expect(response.status).toBe(200);
  expect(text).toContain("Sitemap: https://community.example/sitemap.xml");
  expect(text).toContain("\n\nUser-agent: GPTBot\nDisallow: /\n");
  expect(text).toContain("\n\nUser-agent: CCBot\nDisallow: /\n");
  expect(text).not.toContain("User-agent: OAI-SearchBot");
  expect(text).not.toContain("User-agent: Claude-SearchBot");
  expect(text).toContain("Disallow: /login$\nDisallow: /register$\n");
});

it("honours separate operator crawl preferences", async () => {
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
    search_indexing: false,
    training_crawling: true,
  });
  const text = await (await GET()).text();
  expect(text).toContain("User-agent: *\nDisallow: /");
  expect(text).not.toContain("User-agent: GPTBot");
  expect(text).not.toContain("Sitemap:");
});

it("allows the public site while blocking sections guests cannot list", async () => {
  vi.mocked(getIndexPage).mockResolvedValue(null);
  const text = await (await GET()).text();

  expect(text).toContain("User-agent: *\nAllow: /\n");
  expect(text).toContain("Disallow: /t/\nDisallow: /d/\nDisallow: /l/\n");
  expect(text).not.toContain("Sitemap:");
});

it.each([
  ["threads", "Disallow: /t/\nDisallow: /d/", "Disallow: /l/"],
  ["library", "Disallow: /l/", "Disallow: /t/"],
])(
  "only blocks the inaccessible %s section",
  async (kind, blocked, accessible) => {
    vi.mocked(getIndexPage).mockImplementation(async (requestedKind) => {
      if (requestedKind === kind) {
        return null;
      }

      return { entries: [], totalPages: 1 };
    });
    const text = await (await GET()).text();

    expect(text).toContain("User-agent: *\nAllow: /\n");
    expect(text).toContain(blocked);
    expect(text).not.toContain(accessible);
    expect(text).toContain("Sitemap: https://community.example/sitemap.xml");
  },
);

it("blocks crawling without querying content when search indexing is disabled", async () => {
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
    search_indexing: false,
    training_crawling: true,
  });
  vi.mocked(getIndexPage).mockRejectedValue(new Error("unavailable"));
  const response = await GET();

  expect(response.status).toBe(200);
  expect(await response.text()).toBe("User-agent: *\nDisallow: /\n");
  expect(getIndexPage).not.toHaveBeenCalled();
});

it("returns an error when the API is unavailable", async () => {
  vi.mocked(getIndexPage).mockRejectedValue(new Error("unavailable"));
  vi.spyOn(console, "error").mockImplementation(() => {});
  expect((await GET()).status).toBe(503);
});

it("serves custom robots.txt verbatim without querying content", async () => {
  const robots = "User-agent: *\nDisallow: /\n";
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
    robots_txt: robots,
  });
  const response = await GET();
  expect(await response.text()).toBe(robots);
  expect(getIndexPage).not.toHaveBeenCalled();
});

it("does not fall back to permissive defaults when settings cannot be read", async () => {
  vi.mocked(getDiscoverySettings).mockRejectedValue(new Error("unavailable"));
  vi.spyOn(console, "error").mockImplementation(() => {});
  expect((await GET()).status).toBe(503);
});

it("previews the generated policy without replacing a saved custom document", async () => {
  const customRobots = "User-agent: *\nDisallow: /custom-only/\n";
  vi.mocked(getDiscoverySettings).mockResolvedValue({
    ...DefaultDiscoverySettings,
    search_indexing: false,
    training_crawling: true,
    robots_txt: customRobots,
  });

  const preview = await generateRobotsTxt(false, true);
  expect(preview).toBe("User-agent: *\nDisallow: /\n");

  const served = await GET();
  expect(await served.text()).toBe(customRobots);
});
