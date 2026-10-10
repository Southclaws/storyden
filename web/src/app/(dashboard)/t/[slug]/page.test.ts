import assert from "node:assert/strict";
import { afterEach, test, vi } from "vitest";

vi.mock("@/api/openapi-server/threads", () => ({ threadGet: vi.fn() }));
vi.mock("@/auth/server-session", () => ({ getServerSession: vi.fn() }));
vi.mock("@/lib/settings/settings-server", () => ({
  getSettings: async () => ({ title: "Community" }),
}));
vi.mock("@/screens/thread/ThreadScreen/ThreadScreen", () => ({
  ThreadScreen: () => null,
}));

import { generateMetadata } from "./page";

const props = (page?: string) => ({
  params: Promise.resolve({ slug: "first-post" }),
  searchParams: Promise.resolve(page ? { page } : {}),
});

afterEach(() => vi.unstubAllGlobals());

test("a published thread with zero replies stays indexable", async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          slug: "first-post",
          title: "First post",
          description: "Welcome",
          visibility: "published",
          createdAt: "2026-01-02T03:04:05Z",
          updatedAt: "2026-01-03T03:04:05Z",
          replies: { total_pages: 0 },
        }),
        { status: 200 },
      ),
    ),
  );

  const metadata = await generateMetadata(props());
  assert.equal(metadata.title, "First post | Community");
  assert.deepEqual(
    metadata.alternates?.canonical,
    "http://localhost:3000/t/first-post",
  );
  assert.equal(metadata.robots, undefined);
  assert.equal(metadata.pagination?.next, undefined);

  const missingPage = await generateMetadata(props("2"));
  assert.deepEqual(missingPage.robots, { index: false, follow: false });
});

test("unpublished and guest-forbidden threads expose no content metadata", async () => {
  const fetch = vi
    .fn()
    .mockResolvedValueOnce(
      new Response(
        JSON.stringify({ visibility: "draft", title: "Secret draft" }),
        { status: 200 },
      ),
    )
    .mockResolvedValueOnce(new Response(null, { status: 403 }));
  vi.stubGlobal("fetch", fetch);

  for (const metadata of [
    await generateMetadata(props()),
    await generateMetadata(props()),
  ]) {
    assert.deepEqual(metadata.robots, { index: false, follow: false });
    assert.equal(JSON.stringify(metadata).includes("Secret draft"), false);
    assert.equal(metadata.alternates, undefined);
  }
});
