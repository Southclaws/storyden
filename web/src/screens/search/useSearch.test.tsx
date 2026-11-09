import { act, renderHook, waitFor } from "@testing-library/react";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { type PropsWithChildren } from "react";
import { SWRConfig } from "swr";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetcher } from "@/api/client";

import { useSearchScreen } from "./useSearch";

vi.mock("@/api/client", () => ({ fetcher: vi.fn() }));

const initialCategories = [
  { name: "General", slug: "general", id: "general-id" },
  { name: "Help", slug: "questions", id: "help-id" },
];

function Wrapper({ children }: PropsWithChildren) {
  return (
    <NuqsTestingAdapter hasMemory>
      <SWRConfig
        value={{
          provider: () => new Map(),
          errorRetryCount: 0,
          revalidateOnFocus: false,
        }}
      >
        {children}
      </SWRConfig>
    </NuqsTestingAdapter>
  );
}

function renderSearch() {
  return renderHook(
    () =>
      useSearchScreen({
        initialQuery: "",
        initialPage: 1,
        initialKind: [],
        initialAuthors: [],
        initialCategories: [],
        initialTags: [],
      }),
    { wrapper: Wrapper },
  );
}

describe("search filter options", () => {
  beforeEach(() => {
    vi.mocked(fetcher).mockImplementation(async (url) => {
      const { pathname, searchParams } = new URL(url, "http://localhost");
      const query = searchParams.get("q");

      if (pathname === "/profiles") {
        return { profiles: [{ handle: query || "alice" }] };
      }

      if (pathname === "/tags") {
        return { tags: [{ name: query || "gardening" }] };
      }

      if (pathname === "/categories") {
        return { categories: initialCategories };
      }

      throw new Error(`Unexpected request: ${url}`);
    });
  });

  it("loads options for every filter before the user types", async () => {
    const { result } = renderSearch();

    await waitFor(() => {
      expect(result.current.filters.authorsResults).toEqual([
        { label: "alice", value: "alice" },
      ]);
      expect(result.current.filters.categoriesResults).toEqual([
        { label: "General", value: "general" },
        { label: "Help", value: "questions" },
      ]);
      expect(result.current.filters.tagsResults).toEqual([
        { label: "gardening", value: "gardening" },
      ]);
    });
  });

  it("filters categories by name or slug and restores them on clearing", async () => {
    const { result } = renderSearch();
    await waitFor(() =>
      expect(result.current.filters.categoriesResults).toHaveLength(2),
    );

    for (const query of ["HELP", "questions"]) {
      act(() => result.current.handlers.handleQueryCategories(query));
      expect(result.current.filters.categoriesResults).toEqual([
        { label: "Help", value: "questions" },
      ]);
    }

    act(() => result.current.handlers.handleQueryCategories("missing"));
    expect(result.current.filters.categoriesResults).toEqual([]);

    act(() => result.current.handlers.handleQueryCategories(""));
    expect(result.current.filters.categoriesResults).toHaveLength(2);
  });

  it("applies a category query even when options arrive after typing", async () => {
    let resolveCategories!: (value: unknown) => void;
    vi.mocked(fetcher).mockReturnValue(
      new Promise((resolve) => {
        resolveCategories = resolve;
      }),
    );
    const { result } = renderSearch();

    expect(result.current.filters.categoriesLoading).toBe(true);
    act(() => result.current.handlers.handleQueryCategories("help"));
    await act(async () => {
      resolveCategories({ categories: initialCategories });
    });

    expect(result.current.filters.categoriesResults).toEqual([
      { label: "Help", value: "questions" },
    ]);
  });

  it.each([
    {
      path: "/profiles",
      handler: "handleQueryAuthors",
      field: "authorsResults",
      collection: "profiles",
      item: { handle: "late" },
      initial: "alice",
    },
    {
      path: "/tags",
      handler: "handleQueryTags",
      field: "tagsResults",
      collection: "tags",
      item: { name: "late" },
      initial: "gardening",
    },
  ] as const)(
    "restores $path options without a late search replacing them",
    async ({ path, handler, field, collection, item, initial }) => {
      const { result } = renderSearch();
      const initialOptions = [{ label: initial, value: initial }];
      await waitFor(() =>
        expect(result.current.filters[field]).toEqual(initialOptions),
      );

      let resolveSearch!: (value: unknown) => void;
      vi.mocked(fetcher).mockImplementationOnce(
        () =>
          new Promise((resolve) => {
            resolveSearch = resolve;
          }),
      );
      act(() => result.current.handlers[handler]("late"));
      await waitFor(() =>
        expect(fetcher).toHaveBeenCalledWith(`${path}?q=late`, {
          method: "GET",
        }),
      );

      act(() => result.current.handlers[handler](""));
      expect(result.current.filters[field]).toEqual(initialOptions);
      await act(async () => {
        resolveSearch({ [collection]: [item] });
      });
      expect(result.current.filters[field]).toEqual(initialOptions);
    },
  );

  it("cancels queued searches when inputs are cleared", async () => {
    const { result } = renderSearch();
    await waitFor(() =>
      expect(result.current.filters.tagsResults).toHaveLength(1),
    );

    act(() => {
      result.current.handlers.handleQueryAuthors("cancelled");
      result.current.handlers.handleQueryTags("cancelled");
      result.current.handlers.handleQueryAuthors("");
      result.current.handlers.handleQueryTags("");
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 350));
    });

    expect(
      vi.mocked(fetcher).mock.calls.some(([url]) => url.includes("cancelled")),
    ).toBe(false);
    expect(result.current.filters.authorsResults).toEqual([
      { label: "alice", value: "alice" },
    ]);
    expect(result.current.filters.tagsResults).toEqual([
      { label: "gardening", value: "gardening" },
    ]);
  });
});
