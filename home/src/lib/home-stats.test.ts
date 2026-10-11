import assert from "node:assert/strict";
import { afterEach, mock, test } from "node:test";

import { getHomeStats } from "./home-stats";
import snapshot from "./home-stats.snapshot.json";

afterEach(() => mock.restoreAll());

test("counts commits and contributors from last-page links, not issues or contribution sums", async () => {
  mock.method(globalThis, "fetch", async (input: string) => {
    const url = new URL(input);
    if (url.pathname.endsWith("/commits")) {
      assert.equal(url.searchParams.get("per_page"), "1");
      return Response.json([{ sha: "abc" }], {
        headers: {
          link: '<https://api.github.com/repositories/1/commits?per_page=1&page=4201>; rel="last", <https://api.github.com/repositories/1/commits?per_page=1&page=2>; rel="next"',
        },
      });
    }
    if (url.pathname.endsWith("/contributors")) {
      assert.equal(url.searchParams.get("per_page"), "1");
      return Response.json([{ contributions: 700 }], {
        headers: {
          link: '<https://api.github.com/repositories/1/contributors?per_page=1&page=105>; rel="last"',
        },
      });
    }
    return Response.json({ stargazers_count: 350, open_issues_count: 99 });
  });
  const result = await getHomeStats();
  assert.equal(result.stars, 350);
  assert.equal(result.commits, 4201);
  assert.equal(result.contributors, 105);
});

test("uses list length when GitHub returns a single page", async () => {
  mock.method(globalThis, "fetch", async (input: string) => {
    if (input.includes("/commits?")) return Response.json([{ sha: "abc" }]);
    if (input.includes("/contributors?")) return Response.json([]);
    return Response.json({ stargazers_count: 0 });
  });
  const result = await getHomeStats();
  assert.equal(result.stars, 0);
  assert.equal(result.commits, 1);
  assert.equal(result.contributors, 0);
});

test("a rate-limited metric retains its measurement without discarding successful requests", async () => {
  mock.method(console, "warn", () => {});
  mock.method(globalThis, "fetch", async (input: string) => {
    if (input.includes("/commits?")) return new Response(null, { status: 403 });
    if (input.includes("/contributors?")) return Response.json([{}]);
    return Response.json({ stargazers_count: 351 });
  });
  const result = await getHomeStats();
  assert.equal(result.stars, 351);
  assert.equal(result.commits, snapshot.commits);
  assert.equal(result.contributors, 1);
  assert.ok(result.apis > 0);
  assert.ok(result.loc > 0);
  await assert.rejects(getHomeStats({ strict: true }), /403/);
});

test("malformed counts and network failures never replace measured values", async () => {
  mock.method(console, "warn", () => {});
  mock.method(globalThis, "fetch", async (input: string) => {
    if (input.includes("/commits?")) throw new Error("timeout");
    if (input.includes("/contributors?")) {
      return Response.json([{}], {
        headers: {
          link: '<https://api.github.com/repositories/1/contributors?page=oops>; rel="last"',
        },
      });
    }
    return Response.json({ stargazers_count: "invalid" });
  });
  const result = await getHomeStats();
  assert.equal(result.stars, snapshot.stars);
  assert.equal(result.commits, snapshot.commits);
  assert.equal(result.contributors, snapshot.contributors);
});
