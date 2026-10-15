import type { ReactNode } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { beforeEach, expect, it, vi } from "vitest";

import CategoryPage from "@/app/(dashboard)/d/[category]/page";
import LibraryPage from "@/app/(dashboard)/l/[...slug]/page";
import LibraryIndexPage from "@/app/(dashboard)/l/page";
import ProfilePage from "@/app/(dashboard)/m/[handle]/page";
import ThreadPage, {
  generateMetadata as threadMetadata,
} from "@/app/(dashboard)/t/[slug]/page";
import TagPage from "@/app/(dashboard)/tags/[tag]/page";

vi.mock("@/components/content/ContentRenderer/ServerRenderedContent", () => ({
  ServerRenderedContent: ({ children }: { children: ReactNode }) => children,
}));
vi.mock("@/config", () => ({
  API_ADDRESS: "https://api.example",
  WEB_ADDRESS: "https://community.example",
  getAPIAddress: () => "https://api.example",
}));
vi.mock("next/headers", () => ({
  cookies: async () => ({
    getAll: () => [{ name: "session", value: "admin" }],
  }),
  headers: async () => new Headers({ "x-forwarded-for": "203.0.113.1" }),
}));
vi.mock("next/navigation", () => ({
  notFound: () => {
    throw new Error("not found");
  },
}));
vi.mock("@/auth/server-session", () => ({
  getServerSession: async () => undefined,
}));
vi.mock("@/lib/settings/settings-server", () => ({
  getSettings: async () => ({
    title: "Community",
    metadata: { signatures: {} },
    discovery: { search_indexing: false },
  }),
}));
vi.mock("@/components/member/MemberBadge/MemberAvatarEditable", () => ({
  MemberAvatarEditable: () => null,
}));
vi.mock("@/components/site/Unready", () => ({ UnreadyBanner: () => null }));
vi.mock("@/screens/category/CategoryScreen", () => ({
  CategoryScreen: () => null,
}));
vi.mock("@/screens/library/LibraryIndexScreen", () => ({
  LibraryIndexScreen: () => null,
}));
vi.mock("@/screens/library/LibraryPageScreen/LibraryPageScreen", () => ({
  LibraryPageScreen: () => null,
}));
vi.mock("@/screens/profile/ProfileScreen", () => ({
  ProfileScreen: () => null,
}));
vi.mock("@/screens/tags/TagScreen", () => ({ TagScreen: () => null }));
vi.mock("@/screens/thread/ThreadScreen/ThreadScreen", () => ({
  ThreadScreen: () => null,
}));

const author = { id: "author", name: "Ada", handle: "ada" };
const dates = {
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-02-01T00:00:00Z",
};
const profile = {
  ...author,
  ...dates,
  joined: dates.createdAt,
  updatedAt: "0001-01-01T00:00:00Z",
  bio: "Public biography",
};
const node = {
  ...dates,
  id: "node",
  name: "Drawing & colour",
  slug: "drawing-colour",
  description: "A field guide",
  visibility: "published",
  hide_child_tree: false,
  ancestors: [{ id: "parent", name: "Art", slug: "art" }],
  children: [],
};
const reply = {
  ...dates,
  id: "reply-two",
  body: "The second page reply",
  author,
  visibility: "published",
};
const thread = {
  ...dates,
  id: "thread",
  slug: "thread-a-conversation",
  title: "A conversation",
  body: "The original post",
  description: "A public discussion",
  visibility: "published",
  author,
  replies: { current_page: 2, total_pages: 2, results: 40, replies: [reply] },
};

let responses: Record<string, { signed: unknown; public: unknown }>;
let publicStatus: number;
let requests: Request[];

beforeEach(() => {
  publicStatus = 200;
  requests = [];
  responses = {
    "/profiles/ada": {
      signed: { ...profile, bio: "Private biography" },
      public: profile,
    },
    "/nodes/drawing-colour": {
      signed: {
        ...node,
        ancestors: [{ id: "private", name: "Private parent", slug: "private" }],
      },
      public: node,
    },
    "/nodes": { signed: { nodes: [node] }, public: { nodes: [node] } },
    "/categories/art": {
      signed: { name: "Art", slug: "art" },
      public: { name: "Art", slug: "art", description: "Creative discussions" },
    },
    "/threads": {
      signed: { threads: [thread] },
      public: { threads: [thread], total_pages: 2 },
    },
    "/threads/thread-a-conversation": { signed: thread, public: thread },
    "/tags/design": {
      signed: {
        name: "design",
        items: [
          { kind: "node", ref: node },
          { kind: "thread", ref: thread },
        ],
      },
      public: {
        name: "design",
        items: [
          { kind: "node", ref: node },
          { kind: "thread", ref: thread },
        ],
      },
    },
  };

  // Exercise generated clients and the real server fetcher. Only the network
  // and unrelated screen rendering are replaced at the test boundary.
  vi.spyOn(globalThis, "fetch").mockImplementation(async (input) => {
    const request = input as Request;
    requests.push(request);
    const path = new URL(request.url).pathname.replace(/^\/api/, "");
    const response = responses[path];

    if (!response) {
      throw new Error(`Unexpected API request: ${path}`);
    }

    if (request.headers.has("Cookie")) {
      return Response.json(response.signed);
    }

    expect(request.headers.has("Authorization")).toBe(false);
    expect(request.headers.has("x-forwarded-for")).toBe(false);
    expect(request.cache).toBe("no-store");

    return Response.json(
      publicStatus === 200 ? response.public : { title: "Unavailable" },
      { status: publicStatus },
    );
  });
});

function structuredData(element: ReactNode) {
  const html = renderToStaticMarkup(element);
  return [
    ...html.matchAll(/<script type="application\/ld\+json">(.*?)<\/script>/gs),
  ].flatMap((match) => JSON.parse(match[1]!));
}

it.each([
  { page: "3", totalPages: 2 },
  { page: "999", totalPages: 2 },
  { page: "2", totalPages: 0 },
])(
  "omits thread structured data for page $page with $totalPages total pages",
  async ({ page, totalPages }) => {
    const emptyPage = {
      ...thread,
      replies: {
        current_page: Number(page),
        total_pages: totalPages,
        results: 0,
        replies: [],
      },
    };
    responses["/threads/thread-a-conversation"] = {
      signed: emptyPage,
      public: emptyPage,
    };
    const props = {
      params: Promise.resolve({ slug: thread.slug }),
      searchParams: Promise.resolve({ page }),
    };

    expect(structuredData(await ThreadPage(props))).toEqual([]);
    expect(await threadMetadata(props)).toEqual({
      title: "Thread Not Found",
      description: "The thread you are looking for does not exist.",
      robots: { index: false, follow: false },
    });
  },
);

it("retains structured data on the first page of a thread without replies", async () => {
  const emptyThread = {
    ...thread,
    replies: { current_page: 1, total_pages: 0, results: 0, replies: [] },
  };
  responses["/threads/thread-a-conversation"] = {
    signed: emptyThread,
    public: emptyThread,
  };
  const props = {
    params: Promise.resolve({ slug: thread.slug }),
    searchParams: Promise.resolve({}),
  };

  expect(structuredData(await ThreadPage(props))).toEqual([
    expect.objectContaining({
      "@type": "DiscussionForumPosting",
      commentCount: 0,
      comment: [],
    }),
  ]);
  expect((await threadMetadata(props)).alternates?.canonical).toBe(
    "https://community.example/t/thread-a-conversation",
  );
});

it("describes the anonymous profile, even when search indexing is disabled", async () => {
  const element = await ProfilePage({
    params: Promise.resolve({ handle: "ada" }),
  });

  expect(structuredData(element)).toEqual([
    {
      "@context": "https://schema.org",
      "@type": "ProfilePage",
      url: "https://community.example/m/ada",
      dateCreated: "2026-01-01T00:00:00Z",
      mainEntity: {
        "@type": "Person",
        name: "Ada",
        alternateName: "ada",
        identifier: "author",
        url: "https://community.example/m/ada",
        description: "Public biography",
        image: "https://api.example/api/accounts/ada/avatar",
      },
    },
  ]);
});

it("builds breadcrumbs from anonymous ancestry, not URL segments or signed-in ancestors", async () => {
  const element = await LibraryPage({
    params: Promise.resolve({ slug: ["invented-parent", "drawing-colour"] }),
  });
  const schemas = structuredData(element);

  expect(schemas[0]["@type"]).toBe("WebPage");
  expect(schemas[1]).toEqual({
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      {
        "@type": "ListItem",
        position: 1,
        name: "Library",
        item: "https://community.example/l",
      },
      {
        "@type": "ListItem",
        position: 2,
        name: "Art",
        item: "https://community.example/l/art",
      },
      {
        "@type": "ListItem",
        position: 3,
        name: "Drawing & colour",
        item: "https://community.example/l/drawing-colour",
      },
    ],
  });
});

it("describes only the displayed public reply page and retains the total comment count", async () => {
  const draft = {
    ...reply,
    id: "draft",
    body: "Private reply",
    visibility: "draft",
  };
  const notDisplayed = {
    ...reply,
    id: "elsewhere",
    body: "Not on the displayed page",
  };
  responses["/threads/thread-a-conversation"] = {
    signed: {
      ...thread,
      replies: { ...thread.replies, replies: [reply, draft] },
    },
    public: {
      ...thread,
      replies: { ...thread.replies, replies: [reply, draft, notDisplayed] },
    },
  };

  const element = await ThreadPage({
    params: Promise.resolve({ slug: thread.slug }),
    searchParams: Promise.resolve({ page: "2" }),
  });
  const [schema] = structuredData(element);

  expect(schema.author).toEqual({
    "@type": "Person",
    name: "Ada",
    url: "https://community.example/m/ada",
  });
  expect(schema.commentCount).toBe(40);
  expect(schema.comment).toEqual([
    {
      "@type": "Comment",
      text: "The second page reply",
      url: "https://community.example/t/thread-a-conversation?page=2#reply-two",
      datePublished: "2026-01-01T00:00:00Z",
      dateModified: "2026-02-01T00:00:00Z",
      author: {
        "@type": "Person",
        name: "Ada",
        url: "https://community.example/m/ada",
      },
    },
  ]);
  expect(
    requests.filter((request) => !request.headers.has("Cookie")),
  ).toHaveLength(1);
  expect(new URL(requests[1]!.url).searchParams.get("page")).toBe("2");
});

it("describes the current category page with a bounded public item list", async () => {
  responses["/threads"]!.public = {
    threads: [thread, { ...thread, visibility: "draft" }],
    total_pages: 2,
  };
  const element = await CategoryPage({
    params: Promise.resolve({ category: "art" }),
    searchParams: Promise.resolve({ page: "2" }),
  });

  expect(structuredData(element)).toEqual([
    {
      "@context": "https://schema.org",
      "@type": "CollectionPage",
      name: "Art",
      description: "Creative discussions",
      url: "https://community.example/d/art?page=2",
      mainEntity: {
        "@type": "ItemList",
        itemListElement: [
          {
            "@type": "ListItem",
            position: 1,
            name: "A conversation",
            url: "https://community.example/t/thread-a-conversation",
          },
        ],
      },
    },
  ]);
  const request = requests.find(
    (request) =>
      !request.headers.has("Cookie") &&
      new URL(request.url).pathname.endsWith("/threads"),
  )!;
  expect(new URL(request.url).searchParams.get("visibility")).toBe("published");
  expect(new URL(request.url).searchParams.get("page")).toBe("2");
});

it("describes library roots without flattening descendants or publishing private nodes", async () => {
  responses["/nodes"]!.public = {
    nodes: [
      { ...node, children: [{ ...node, id: "child" }] },
      { ...node, visibility: "draft" },
    ],
  };
  const [schema] = structuredData(await LibraryIndexPage());

  expect(schema["@type"]).toBe("CollectionPage");
  expect(schema.url).toBe("https://community.example/l");
  expect(schema.mainEntity.itemListElement).toEqual([
    {
      "@type": "ListItem",
      position: 1,
      name: "Drawing & colour",
      url: "https://community.example/l/drawing-colour",
    },
  ]);
  expect(new URL(requests[1]!.url).searchParams.get("depth")).toBe("0");
});

it("describes each tag item with its actual public route", async () => {
  responses["/tags/design"]!.public = {
    name: "design & café",
    items: [
      { kind: "node", ref: node },
      { kind: "thread", ref: thread },
      { kind: "reply", ref: { ...reply, title: "A conversation" } },
      { kind: "profile", ref: profile },
      { kind: "post", ref: thread },
      { kind: "node", ref: { ...node, visibility: "draft" } },
    ],
  };
  responses["/tags/design"]!.signed = responses["/tags/design"]!.public;
  const [schema] = structuredData(
    await TagPage({ params: Promise.resolve({ tag: "design" }) }),
  );

  expect(schema.url).toBe(
    "https://community.example/tags/design%20%26%20caf%C3%A9",
  );
  expect(
    schema.mainEntity.itemListElement.map(
      (item: { position: number; url: string }) => [item.position, item.url],
    ),
  ).toEqual([
    [1, "https://community.example/l/drawing-colour"],
    [2, "https://community.example/t/thread-a-conversation"],
    [3, "https://community.example/t/locate/reply-two"],
    [4, "https://community.example/m/ada"],
    [5, "https://community.example/t/thread-a-conversation"],
  ]);
});

it.each([401, 403, 404, 500])(
  "omits all public schemas when anonymous reads return %i",
  async (status) => {
    publicStatus = status;
    const pages = [
      await ProfilePage({ params: Promise.resolve({ handle: "ada" }) }),
      await LibraryPage({
        params: Promise.resolve({ slug: ["drawing-colour"] }),
      }),
      await ThreadPage({
        params: Promise.resolve({ slug: thread.slug }),
        searchParams: Promise.resolve({}),
      }),
      await CategoryPage({
        params: Promise.resolve({ category: "art" }),
        searchParams: Promise.resolve({}),
      }),
      await LibraryIndexPage(),
      await TagPage({ params: Promise.resolve({ tag: "design" }) }),
    ];

    for (const page of pages) {
      expect(structuredData(page)).toEqual([]);
    }
  },
);

it("does not emit a category collection for an invalid page", async () => {
  expect(
    structuredData(
      await CategoryPage({
        params: Promise.resolve({ category: "art" }),
        searchParams: Promise.resolve({ page: "999" }),
      }),
    ),
  ).toEqual([]);
  expect(
    structuredData(
      await CategoryPage({
        params: Promise.resolve({ category: "art" }),
        searchParams: Promise.resolve({ page: "2junk" }),
      }),
    ),
  ).toEqual([]);
});

it("keeps profile HTML inside one inert JSON-LD script", async () => {
  responses["/profiles/ada"]!.public = {
    ...profile,
    bio: '</script><script>alert("x")</script>',
  };
  const element = await ProfilePage({
    params: Promise.resolve({ handle: "ada" }),
  });

  expect(renderToStaticMarkup(element).match(/<script/g)).toHaveLength(1);
  expect(structuredData(element)[0].mainEntity.description).toBe(
    '</script><script>alert("x")</script>',
  );
});
