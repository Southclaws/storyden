import assert from "node:assert/strict";
import { afterEach, test, vi } from "vitest";

import type { Node, Thread } from "@/api/openapi-schema";

import {
  canonicalPage,
  getAnonymousPublicResource,
  jsonLD,
  libraryPageSchema,
  pageLinks,
  pageNumber,
  threadSchema,
  websiteSchema,
} from "./public";

afterEach(() => vi.unstubAllGlobals());

test("canonical and pagination links use the configured public origin", () => {
  assert.equal(canonicalPage("/t/hello", "1"), "http://localhost:3000/t/hello");
  assert.equal(
    canonicalPage("/t/hello", "2"),
    "http://localhost:3000/t/hello?page=2",
  );
  assert.equal(pageNumber("2junk"), undefined);
  assert.deepEqual(pageLinks("/t/hello", 2, 3), {
    previous: "http://localhost:3000/t/hello",
    next: "http://localhost:3000/t/hello?page=3",
  });
});

test("anonymous metadata excludes drafts and never sends credentials", async () => {
  const fetch = vi
    .fn()
    .mockResolvedValueOnce(
      new Response(JSON.stringify({ visibility: "draft", title: "Secret" }), {
        status: 200,
      }),
    )
    .mockResolvedValueOnce(
      new Response(
        JSON.stringify({ visibility: "published", title: "Public" }),
        {
          status: 200,
        },
      ),
    )
    .mockResolvedValueOnce(new Response(null, { status: 403 }));
  vi.stubGlobal("fetch", fetch);

  assert.equal(await getAnonymousPublicResource("/threads/secret"), undefined);
  assert.deepEqual(await getAnonymousPublicResource("/threads/public"), {
    visibility: "published",
    title: "Public",
  });
  assert.equal(await getAnonymousPublicResource("/threads/closed"), undefined);

  for (const [url, init] of fetch.mock.calls) {
    assert.equal((url as URL).origin, "http://localhost:8000");
    assert.equal((init as RequestInit).cache, "no-store");
    assert.equal((init as RequestInit).credentials, "omit");
    assert.equal((init as RequestInit).headers, undefined);
  }
});

test("JSON-LD represents actual content and cannot escape its script tag", () => {
  const thread = {
    slug: "hello",
    title: "A </script><script>alert(1)</script> post",
    description: "A forum discussion",
    author: { name: "Ada" },
    createdAt: "2026-01-02T03:04:05Z",
    updatedAt: "2026-01-03T03:04:05Z",
  } as Thread;
  const encoded = jsonLD(threadSchema(thread));
  assert.equal(encoded.includes("</script>"), false);
  assert.deepEqual(JSON.parse(encoded), {
    "@context": "https://schema.org",
    "@type": "DiscussionForumPosting",
    headline: thread.title,
    description: thread.description,
    url: "http://localhost:3000/t/hello",
    mainEntityOfPage: "http://localhost:3000/t/hello",
    datePublished: thread.createdAt,
    dateModified: thread.updatedAt,
    author: { "@type": "Person", name: "Ada" },
  });

  const node = {
    slug: "guide",
    name: "Guide",
    description: "How to start",
    createdAt: "2026-01-02T03:04:05Z",
    updatedAt: "2026-01-03T03:04:05Z",
  } as Node;
  assert.equal(libraryPageSchema(node)["@type"], "WebPage");
  assert.equal(websiteSchema("Community", "A place")["@type"], "WebSite");
});
