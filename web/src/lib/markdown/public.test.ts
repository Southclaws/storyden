import assert from "node:assert/strict";
import { afterEach, test, vi } from "vitest";

vi.mock("server-only", () => ({}));
vi.mock("@/config", () => ({ getAPIAddress: () => "http://backend:8000" }));

import { GET as index } from "@/app/index.md/route";
import { GET as library } from "@/app/markdown/l/[...slug]/route";
import { GET as thread } from "@/app/markdown/t/[slug]/route";

import { htmlToPublicMarkdown } from "./public";

afterEach(() => vi.unstubAllGlobals());

const info = {
  title: "A community",
  description: "A place for discussion",
  content: "<p>Welcome <strong>everyone</strong>.</p>",
  web_address: "https://community.example",
  api_address: "https://api.community.example",
  metadata: {
    feed: {
      blocks: [
        { type: "content" },
        { type: "threads", source: "all" },
        { type: "library", layout: "list" },
      ],
    },
  },
};

function routeData(path: string, data: unknown, status = 200) {
  return (url: URL) =>
    url.pathname === path
      ? Response.json(data, { status })
      : url.pathname === "/api/info"
        ? Response.json(info)
        : new Response("missing", { status: 404 });
}

test("rich community content retains paragraphs, code, images and safe links", () => {
  const markdown = htmlToPublicMarkdown(
    '<p>Hello <strong>world</strong>.</p><pre><code>const x = 1;</code></pre><p><a href="https://example.com">source</a> <img src="/image.png" alt="Photo"></p><script>secret()</script><a href="javascript:alert(1)">bad</a>',
  );
  assert.match(markdown, /Hello \*\*world\*\*\./);
  assert.match(markdown, /```\nconst x = 1;/);
  assert.match(markdown, /\[source\]\(https:\/\/example\.com\)/);
  assert.match(markdown, /!\[Photo\]\(\/image\.png\)/);
  assert.doesNotMatch(markdown, /secret\(\)|javascript:/);
  const assets = htmlToPublicMarkdown(
    '<p><img src="/api/assets/photo.png" alt="Photo"> <a href="sdr:node/abc123">page</a></p>',
    {
      web: info.web_address,
      api: info.api_address,
    },
  );
  assert.match(
    assets,
    /https:\/\/api\.community\.example\/api\/assets\/photo\.png/,
  );
  assert.match(
    assets,
    /https:\/\/community\.example\/_\/resolve\/node\/abc123/,
  );
  const relative = htmlToPublicMarkdown(
    '<a href="related">related</a><table><tr><th>Name</th><th>Value</th></tr><tr><td>One</td><td>Two</td></tr></table>',
    {
      web: info.web_address,
      api: info.api_address,
      base: "https://community.example/l/child",
    },
  );
  assert.match(relative, /https:\/\/community\.example\/l\/related/);
  assert.match(
    relative,
    /\| Name \| Value \|\n\| --- \| --- \|\n\| One \| Two \|/,
  );
});

test("thread markdown is fetched anonymously and paginates public replies", async () => {
  const fetcher = vi.fn(async (url: URL, options: RequestInit) => {
    assert.equal(url.host, "backend:8000");
    assert.equal(options.credentials, "omit");
    assert.equal(new Headers(options.headers).get("Cookie"), null);
    if (url.pathname === "/api/info") return Response.json(info);
    return Response.json({
      title: "A thread",
      slug: "a-thread",
      body: "<p>Main text</p>",
      visibility: "published",
      author: { name: "Alice", handle: "alice" },
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-02T00:00:00Z",
      replies: {
        replies: [
          {
            visibility: "published",
            body: "<p>Public reply</p>",
            author: { name: "Bob", handle: "bob" },
            createdAt: "2026-01-03T00:00:00Z",
          },
          {
            visibility: "draft",
            body: "<p>Secret reply</p>",
            author: { name: "Bob", handle: "bob" },
            createdAt: "2026-01-04T00:00:00Z",
          },
        ],
        total_pages: 3,
        next_page: 3,
      },
    });
  });
  vi.stubGlobal("fetch", fetcher);

  const result = await thread(
    new Request("https://community.example/t/a-thread.md?page=2", {
      headers: {
        Cookie: "storyden-session=secret",
        Authorization: "Bearer secret",
      },
    }),
    { params: Promise.resolve({ slug: "a-thread" }) },
  );
  const markdown = await result.text();
  assert.equal(result.status, 200);
  assert.match(result.headers.get("Link") ?? "", /rel="canonical"/);
  assert.match(markdown, /# A thread/);
  assert.match(markdown, /Main text/);
  assert.match(markdown, /Public reply/);
  assert.doesNotMatch(markdown, /Secret reply/);
  assert.match(
    markdown,
    /Next: https:\/\/community\.example\/t\/a-thread\.md\?page=3/,
  );
  assert.equal(fetcher.mock.calls.length, 2);
});

test("draft and unlisted resources never become public Markdown", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: URL) =>
      routeData("/api/threads/private", {
        visibility: "unlisted",
        replies: { replies: [], total_pages: 1 },
      })(url),
    ),
  );
  const threadResponse = await thread(
    new Request("https://community.example/t/private.md"),
    {
      params: Promise.resolve({ slug: "private" }),
    },
  );
  assert.equal(threadResponse.status, 404);

  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: URL) =>
      routeData("/api/nodes/private", {
        visibility: "draft",
        children: [],
      })(url),
    ),
  );
  const nodeResponse = await library(
    new Request("https://community.example/l/private.md"),
    {
      params: Promise.resolve({ slug: ["private"] }),
    },
  );
  assert.equal(nodeResponse.status, 404);
});

test("a published Library child can be read despite a private ancestor", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: URL) =>
      routeData("/api/nodes/child", {
        name: "Child page",
        slug: "child",
        description: "Visible page",
        visibility: "published",
        content: "<p>Read me</p>",
        owner: { name: "Alice", handle: "alice" },
        createdAt: "2026-01-01T00:00:00Z",
        updatedAt: "2026-01-01T00:00:00Z",
        ancestors: [{ slug: "parent", visibility: "draft" }],
        children: [],
      })(url),
    ),
  );
  const response = await library(
    new Request("https://community.example/l/parent/child.md"),
    {
      params: Promise.resolve({ slug: ["parent", "child"] }),
    },
  );
  assert.equal(response.status, 200);
  const markdown = await response.text();
  assert.match(markdown, /Read me/);
  assert.match(markdown, /Canonical: https:\/\/community\.example\/l\/child/);
  assert.doesNotMatch(markdown, /Canonical: .*parent/);
});

test("hidden-tree Library directory uses paginated public children endpoint", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: URL) => {
      if (url.pathname === "/api/info") return Response.json(info);
      if (url.pathname === "/api/nodes/directory")
        return Response.json({
          name: "Directory",
          slug: "directory",
          description: "",
          visibility: "published",
          owner: { name: "Alice", handle: "alice" },
          createdAt: "2026-01-01",
          updatedAt: "2026-01-01",
          content: "",
          hide_child_tree: true,
          children: [],
          meta: { layout: { blocks: [{ type: "directory" }] } },
        });
      if (url.pathname === "/api/nodes/directory/children")
        return Response.json({
          nodes: [
            { name: "Visible child", slug: "visible", visibility: "published" },
            { name: "Private child", slug: "private", visibility: "draft" },
          ],
          total_pages: 2,
          next_page: 2,
        });
      return new Response("missing", { status: 404 });
    }),
  );
  const response = await library(
    new Request("https://community.example/l/directory.md"),
    {
      params: Promise.resolve({ slug: ["directory"] }),
    },
  );
  const markdown = await response.text();
  assert.equal(response.status, 200);
  assert.match(markdown, /Visible child/);
  assert.match(
    markdown,
    /Next: https:\/\/community\.example\/l\/directory\.md\?page=2/,
  );
  assert.doesNotMatch(markdown, /Private child/);
});

test("upstream access and missing statuses are preserved; outages and invalid pages fail", async () => {
  for (const status of [401, 403, 404, 500]) {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (url: URL) =>
        routeData("/api/threads/missing", {}, status)(url),
      ),
    );
    const response = await thread(
      new Request("https://community.example/t/missing.md"),
      {
        params: Promise.resolve({ slug: "missing" }),
      },
    );
    assert.equal(response.status, status === 500 ? 503 : status);
  }
  const invalid = await thread(
    new Request("https://community.example/t/missing.md?page=1000001"),
    {
      params: Promise.resolve({ slug: "missing" }),
    },
  );
  assert.equal(invalid.status, 400);
});

test("index reflects configured public feed and page navigation", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url: URL) => {
      if (url.pathname === "/api/info") return Response.json(info);
      if (url.pathname === "/api/threads")
        return Response.json({
          threads: [
            {
              title: "Public topic",
              slug: "public",
              visibility: "published",
              author: { name: "Alice" },
              createdAt: "2026-01-01",
            },
            {
              title: "Hidden topic",
              slug: "hidden",
              visibility: "draft",
              author: { name: "Alice" },
              createdAt: "2026-01-01",
            },
          ],
          total_pages: 2,
          next_page: 2,
        });
      if (url.pathname === "/api/nodes")
        return Response.json({
          nodes: [
            {
              name: "Public page",
              slug: "public-page",
              description: "Details",
              visibility: "published",
            },
            {
              name: "Hidden page",
              slug: "hidden-page",
              description: "Private",
              visibility: "draft",
            },
          ],
        });
      return new Response("missing", { status: 404 });
    }),
  );
  const response = await index(
    new Request("https://community.example/index.md"),
  );
  const markdown = await response.text();
  assert.equal(response.status, 200);
  assert.match(markdown, /Welcome \*\*everyone\*\*/);
  assert.match(markdown, /Public topic/);
  assert.match(markdown, /Public page/);
  assert.match(
    markdown,
    /Next: https:\/\/community\.example\/index\.md\?page=2/,
  );
  assert.doesNotMatch(markdown, /Hidden topic|Hidden page/);
  const pastEnd = await index(
    new Request("https://community.example/index.md?page=3"),
  );
  assert.equal(pastEnd.status, 404);
});
