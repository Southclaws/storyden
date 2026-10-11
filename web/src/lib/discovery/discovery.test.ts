import assert from "node:assert/strict";
import { afterEach, test, vi } from "vitest";

vi.mock("server-only", () => ({}));
vi.mock("@/config", () => ({ getAPIAddress: () => "http://backend:8000" }));
vi.mock("next/headers", () => ({
  cookies: async () => ({ getAll: () => [] }),
  headers: async () => new Headers(),
}));

import { GET as authGuide } from "@/app/auth.md/route";
import { GET as llms } from "@/app/llms.txt/route";

import { proxyDiscovery } from "./proxy";
import { getAICatalog } from "./public";

afterEach(() => vi.unstubAllGlobals());

test("catalog proxy forwards only conditional request headers and public response headers", async () => {
  const fetcher = vi.fn(async (_url: URL, options: RequestInit) => {
    assert.equal(options.method, "GET");
    assert.equal(new Headers(options.headers).get("If-None-Match"), '"old"');
    assert.equal(new Headers(options.headers).get("Cookie"), null);
    assert.equal(new Headers(options.headers).get("Authorization"), null);
    return new Response(null, {
      status: 304,
      headers: {
        ETag: '"new"',
        "Cache-Control": "public, max-age=3600",
        "Set-Cookie": "private=yes",
      },
    });
  });
  vi.stubGlobal("fetch", fetcher);

  const response = await proxyDiscovery(
    new Request("https://community.example/.well-known/ai-catalog.json", {
      headers: {
        Cookie: "storyden-session=secret",
        Authorization: "Bearer secret",
        "If-None-Match": '"old"',
      },
    }),
    "/.well-known/ai-catalog.json",
  );

  assert.equal(
    fetcher.mock.calls[0]?.[0].href,
    "http://backend:8000/.well-known/ai-catalog.json",
  );
  assert.equal(response.status, 304);
  assert.equal(await response.text(), "");
  assert.equal(response.headers.get("ETag"), '"new"');
  assert.equal(response.headers.get("Set-Cookie"), null);
});

test("catalog proxy preserves media, CORS, HEAD and upstream failure status", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(
      async () =>
        new Response('{"entries":[]}', {
          headers: {
            "Content-Type": "application/ai-catalog+json",
            "Access-Control-Allow-Origin": "*",
            Link: "<https://community.example/.well-known/api-catalog>; rel=api-catalog",
          },
        }),
    ),
  );
  const get = await proxyDiscovery(
    new Request("https://community.example/x"),
    "/.well-known/ai-catalog.json",
  );
  assert.equal(get.headers.get("Content-Type"), "application/ai-catalog+json");
  assert.equal(get.headers.get("Access-Control-Allow-Origin"), "*");
  assert.match(get.headers.get("Link") ?? "", /api-catalog/);
  assert.equal(await get.text(), '{"entries":[]}');

  const head = await proxyDiscovery(
    new Request("https://community.example/x", { method: "HEAD" }),
    "/.well-known/ai-catalog.json",
  );
  assert.equal(head.status, 200);
  assert.equal(await head.text(), "");

  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response("missing", { status: 404 })),
  );
  const missing = await proxyDiscovery(
    new Request("https://community.example/x"),
    "/.well-known/api-catalog",
  );
  assert.equal(missing.status, 404);

  vi.stubGlobal(
    "fetch",
    vi.fn(async () => {
      throw new Error("down");
    }),
  );
  const unavailable = await proxyDiscovery(
    new Request("https://community.example/x"),
    "/.well-known/api-catalog",
  );
  assert.equal(unavailable.status, 502);
  assert.equal(unavailable.headers.get("Cache-Control"), "no-store");
});

const info = {
  title: "Community",
  description: "A place to meet",
  web_address: "https://community.example",
  api_address: "https://api.community.example",
  capabilities: ["oauth"],
};

function discoveryFetch(entries: unknown, oauth = true) {
  return vi.fn(async (input: Request | URL) => {
    const url = input instanceof Request ? new URL(input.url) : input;
    if (url.pathname === "/api/info")
      return Response.json({
        ...info,
        capabilities: oauth ? ["oauth"] : [],
      });
    if (url.pathname === "/.well-known/ai-catalog.json") {
      return Response.json({ specVersion: "1.0", entries });
    }
    return new Response(null, { status: 404 });
  });
}

test("instance guides use configured origins and advertise MCP only when enabled", async () => {
  vi.stubGlobal(
    "fetch",
    discoveryFetch([
      {
        identifier: "urn:air:community.example:mcp:storyden",
        type: "application/mcp-server-card+json",
        url: "https://api.community.example/mcp/server-card",
      },
    ]),
  );
  const enabled = await (await llms()).text();
  assert.match(enabled, /https:\/\/community\.example\/auth\.md/);
  assert.match(
    enabled,
    /https:\/\/api\.community\.example\/api\/openapi\.json/,
  );
  assert.match(enabled, /https:\/\/api\.community\.example\/mcp\/server-card/);
  assert.doesNotMatch(enabled, /backend:8000/);

  vi.stubGlobal("fetch", discoveryFetch([]));
  assert.doesNotMatch(await (await llms()).text(), /MCP Server Card/);

  vi.stubGlobal(
    "fetch",
    discoveryFetch([
      {
        identifier: "urn:air:community.example:mcp:external",
        type: "application/mcp-server-card+json",
        url: "https://mcp.provider.example/custom-card.json",
      },
    ]),
  );
  assert.match(
    await (await llms()).text(),
    /https:\/\/mcp\.provider\.example\/custom-card\.json/,
  );
});

test("old or unavailable catalogs do not produce false MCP advertisements", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (input: Request | URL) => {
      const url = input instanceof Request ? new URL(input.url) : input;
      return url.pathname === "/api/info"
        ? Response.json(info)
        : new Response("missing", { status: 404 });
    }),
  );
  assert.doesNotMatch(await (await llms()).text(), /MCP Server Card/);

  vi.stubGlobal(
    "fetch",
    vi.fn(async () => {
      throw new Error("down");
    }),
  );
  assert.equal((await llms()).status, 503);
});

test("malformed catalog entries do not produce card links", async () => {
  const entry = {
    identifier: "urn:air:community.example:mcp:storyden",
    type: "application/mcp-server-card+json",
    url: "https://api.community.example/mcp/server-card",
  };
  for (const entries of [
    [null],
    [{ type: entry.type, url: entry.url }],
    [{ ...entry, url: "not a URL" }],
    [{ ...entry, data: {} }],
    [{ identifier: entry.identifier, type: entry.type }],
  ]) {
    vi.stubGlobal("fetch", discoveryFetch(entries));
    assert.equal(await getAICatalog(), null);
    assert.doesNotMatch(await (await llms()).text(), /MCP Server Card/);
  }
  vi.stubGlobal(
    "fetch",
    discoveryFetch([{ ...entry, url: "javascript:alert(1)" }]),
  );
  assert.doesNotMatch(
    await (await llms()).text(),
    /javascript:|MCP Server Card/,
  );
});

test("inline cards are discoverable through the catalog without inventing a card URL", async () => {
  vi.stubGlobal(
    "fetch",
    discoveryFetch([
      {
        identifier: "urn:air:community.example:mcp:storyden",
        type: "application/mcp-server-card+json",
        data: {
          $schema:
            "https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json",
          name: "example.community/mcp",
          description: "Community MCP server",
          version: "1.0.0",
        },
      },
    ]),
  );
  const response = await llms();
  assert.equal(response.status, 200);
  const body = await response.text();
  assert.match(body, /AI catalog/);
  assert.doesNotMatch(body, /MCP Server Card|\/mcp\/server-card|unavailable/);
});

test("catalog consumer accepts compatible minor versions and unknown fields", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async () =>
      Response.json({
        specVersion: "1.1",
        entries: [],
        futureExtension: true,
      }),
    ),
  );
  assert.deepEqual((await getAICatalog())?.entries, []);
  vi.stubGlobal(
    "fetch",
    vi.fn(async () =>
      Response.json({
        specVersion: "2.0",
        entries: [],
      }),
    ),
  );
  assert.equal(await getAICatalog(), null);
});

test("authentication guide reflects the instance OAuth capability", async () => {
  vi.stubGlobal("fetch", discoveryFetch([], false));
  const disabled = await (await authGuide()).text();
  assert.match(disabled, /OAuth is not advertised/);
  assert.doesNotMatch(disabled, /OAuth discovery:/);

  vi.stubGlobal("fetch", discoveryFetch([], true));
  const enabled = await (await authGuide()).text();
  assert.match(
    enabled,
    /api\.community\.example\/\.well-known\/oauth-authorization-server/,
  );
  assert.match(enabled, /Dynamic client registration is controlled/);
});
