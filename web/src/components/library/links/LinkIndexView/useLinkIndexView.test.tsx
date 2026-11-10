import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { SWRConfig } from "swr";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { LinkListResult, LinkReference } from "@/api/openapi-schema";

import { LinkIndexView } from "./LinkIndexView";

const { list, create, onUrlUpdate, auth } = vi.hoisted(() => ({
  list: vi.fn(),
  create: vi.fn(),
  onUrlUpdate: vi.fn(),
  auth: { canUseLinks: true },
}));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("@/auth", () => ({
  useSession: () => ({
    roles: [{ permissions: auth.canUseLinks ? ["USE_LINKS"] : [] }],
  }),
}));

vi.mock("@/api/client", () => ({
  fetcher: (path: string, options: { method: string; body?: string }) => {
    if (options.method === "POST") {
      return create(JSON.parse(options.body ?? "{}"));
    }

    const params = new URL(path, "https://storyden.invalid").searchParams;
    return list(params.get("q") ?? "", Number(params.get("page") ?? 1));
  },
}));

const link: LinkReference = {
  id: "link-a",
  createdAt: "2026-11-09T12:00:00Z",
  updatedAt: "2026-11-09T12:00:00Z",
  url: "https://a.invalid/",
  slug: "link-a",
  domain: "a.invalid",
  title: "A link",
};
const otherLink: LinkReference = {
  ...link,
  id: "link-b",
  url: "https://b.invalid/",
  slug: "link-b",
  domain: "b.invalid",
  title: "B link",
};
const all: LinkListResult = {
  links: [link, otherLink],
  current_page: 1,
  page_size: 2,
  results: 5,
  total_pages: 3,
};
const empty: LinkListResult = { ...all, links: [], results: 0, total_pages: 0 };
const single: LinkListResult = {
  ...all,
  links: [link],
  results: 1,
  total_pages: 1,
};

function renderView(query = "", page = 1, initialResult = all) {
  const searchParams = new URLSearchParams({ q: query, page: page.toString() });

  return render(
    <NuqsTestingAdapter
      hasMemory
      searchParams={searchParams}
      onUrlUpdate={onUrlUpdate}
    >
      <SWRConfig
        value={{
          provider: () => new Map(),
          shouldRetryOnError: false,
          dedupingInterval: 0,
        }}
      >
        <LinkIndexView
          query={query}
          page={page}
          initialResult={initialResult}
        />
      </SWRConfig>
    </NuqsTestingAdapter>,
  );
}

function input() {
  return screen.getByRole("searchbox", { name: "Search links or paste a URL" });
}

function results() {
  return within(screen.getByRole("region", { name: "Link results" }));
}

function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((resolvePromise) => {
    resolve = resolvePromise;
  });

  return { promise, resolve };
}

describe("Links search and add flow", () => {
  beforeEach(() => {
    auth.canUseLinks = true;
    list.mockResolvedValue(all);
    create.mockResolvedValue(link);
  });

  it("searches URLs without offering or submitting creation when Use Links is absent", async () => {
    auth.canUseLinks = false;
    list.mockResolvedValue(empty);
    renderView(link.url, 1, empty);

    expect(
      screen.queryByRole("button", { name: "Add link" }),
    ).not.toBeInTheDocument();
    await userEvent.click(
      screen.getByRole("button", { name: "Search" }),
    );
    fireEvent.submit(input().closest("form")!);

    expect(create).not.toHaveBeenCalled();
    expect(input()).toHaveValue(link.url);
  });

  it("filters while typing and removes pagination that belongs to the previous results", async () => {
    list.mockImplementation(async (query: string) => (query ? single : all));
    renderView();
    expect(results().getByRole("link", { name: "3" })).toBeVisible();

    fireEvent.change(input(), { target: { value: "A link" } });
    expect(
      results().queryByRole("link", { name: "3" }),
    ).not.toBeInTheDocument();
    expect(input()).toHaveValue("A link");

    await waitFor(() => expect(list).toHaveBeenCalledWith("A link", 1));
    await waitFor(() =>
      expect(results().getAllByRole("article")).toHaveLength(1),
    );
    expect(results().getByRole("heading", { name: "A link" })).toBeVisible();
    expect(create).not.toHaveBeenCalled();
  });

  it("searches a pasted existing URL without saving it again", async () => {
    list.mockImplementation(async (query: string) => (query ? single : all));
    const user = userEvent.setup();
    renderView();
    await user.click(input());
    await user.paste(link.url);

    await waitFor(() => expect(list).toHaveBeenCalledWith(link.url, 1));
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled(),
    );
    await user.keyboard("{Enter}");
    expect(create).not.toHaveBeenCalled();
    expect(results().getAllByRole("article")).toHaveLength(1);
  });

  it("shows an empty search for a new URL and waits for explicit Add link", async () => {
    list.mockImplementation(async (query: string) => (query ? empty : all));
    renderView();
    fireEvent.change(input(), { target: { value: "new.invalid" } });

    await waitFor(() =>
      expect(results().getByText("No links were found.")).toBeVisible(),
    );
    expect(screen.getByRole("button", { name: "Add link" })).toBeEnabled();
    expect(create).not.toHaveBeenCalled();
    expect(
      results().queryByRole("link", { name: "2" }),
    ).not.toBeInTheDocument();
  });

  it("shows the saved redirect destination in the search and normal result list", async () => {
    const submitted = "https://www.a.invalid/old";
    list.mockImplementation(async (query: string) =>
      query === link.url ? single : empty,
    );
    renderView(submitted, 1, empty);
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    await waitFor(() => expect(input()).toHaveValue(link.url));
    await waitFor(() =>
      expect(results().getByRole("article", { name: "A link" })).toBeVisible(),
    );
    expect(create).toHaveBeenCalledExactlyOnceWith({ url: submitted });
    expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled();
    expect(results().getAllByRole("article")).toHaveLength(1);
  });

  it("adds immediately, then shows the saved link in the normal filtered list", async () => {
    let saved = false;
    list.mockImplementation(async (query: string) => {
      if (!query) {
        return all;
      }

      return saved ? single : empty;
    });
    create.mockImplementation(async () => {
      saved = true;
      return link;
    });
    renderView();
    fireEvent.change(input(), { target: { value: "a.invalid" } });
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    await waitFor(() =>
      expect(results().getByRole("article", { name: "A link" })).toBeVisible(),
    );
    expect(create).toHaveBeenCalledExactlyOnceWith({ url: link.url });
    expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled();
    expect(results().getAllByRole("article")).toHaveLength(1);
    expect(
      results().queryByRole("link", { name: "3" }),
    ).not.toBeInTheDocument();
  });

  it("prevents repeated submission while adding and lets Clear return to the complete list", async () => {
    const pending = deferred<LinkReference>();
    create.mockReturnValueOnce(pending.promise);
    list.mockImplementation(async (query: string) => (query ? empty : all));
    renderView();
    fireEvent.change(input(), { target: { value: "new.invalid" } });
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));
    expect(screen.getByRole("button", { name: "Adding…" })).toBeDisabled();
    fireEvent.submit(input().closest("form")!);
    expect(create).toHaveBeenCalledTimes(1);

    fireEvent.click(screen.getByRole("button", { name: "Clear search" }));
    expect(input()).toHaveValue("");
    expect(input()).toHaveFocus();
    await act(async () => pending.resolve(link));
    await waitFor(() =>
      expect(results().getAllByRole("article")).toHaveLength(2),
    );
    expect(results().getByRole("link", { name: "3" })).toBeVisible();
    expect(screen.getByRole("button", { name: "Search" })).toBeEnabled();
  });

  it("keeps an add error near the input and retries the same URL", async () => {
    create.mockRejectedValueOnce(new Error("Network unavailable"));
    let saved = false;
    list.mockImplementation(async () => (saved ? single : empty));
    renderView("a.invalid", 1, empty);
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Network unavailable",
    );
    expect(screen.getByRole("button", { name: "Add link" })).toBeEnabled();
    saved = true;
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled(),
    );
    expect(create).toHaveBeenCalledTimes(2);
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("preserves the saved state if refreshing the list fails and offers search recovery", async () => {
    list.mockResolvedValue(empty);
    renderView("a.invalid", 1, empty);
    await waitFor(() => expect(list).toHaveBeenCalled());
    list.mockRejectedValue(new Error("List unavailable"));
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Could not search links",
    );
    expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled();
    expect(create).toHaveBeenCalledTimes(1);

    list.mockResolvedValue(single);
    fireEvent.click(screen.getByRole("button", { name: "Try again" }));
    await waitFor(() =>
      expect(results().getByRole("article", { name: "A link" })).toBeVisible(),
    );
  });

  it("keeps the add operation busy until the saved link is available in the results", async () => {
    const pendingResults = deferred<LinkListResult>();
    list.mockResolvedValue(empty);
    renderView("a.invalid", 1, empty);
    await waitFor(() => expect(list).toHaveBeenCalled());
    list.mockReturnValue(pendingResults.promise);

    fireEvent.click(screen.getByRole("button", { name: "Add link" }));
    await waitFor(() => expect(create).toHaveBeenCalledTimes(1));
    expect(screen.getByRole("button", { name: "Adding…" })).toBeDisabled();
    expect(
      results().queryByText("No links were found."),
    ).not.toBeInTheDocument();

    await act(async () => pendingResults.resolve(single));
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled(),
    );
    expect(results().getByRole("article", { name: "A link" })).toBeVisible();
  });

  it("keeps a newer search when a pending save returns a redirect destination", async () => {
    const pending = deferred<LinkReference>();
    create.mockReturnValueOnce(pending.promise);
    list.mockImplementation(async (query: string) =>
      query === "browser" ? { ...single, links: [otherLink] } : empty,
    );
    renderView("https://www.a.invalid/old", 1, empty);
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));
    fireEvent.change(input(), { target: { value: "browser" } });

    await act(async () => pending.resolve(link));
    await waitFor(() =>
      expect(results().getByRole("article", { name: "B link" })).toBeVisible(),
    );
    expect(input()).toHaveValue("browser");
    expect(onUrlUpdate.mock.lastCall![0].searchParams.get("q")).toBe("browser");
  });

  it("keeps a newer add pending when an earlier request for the same URL completes", async () => {
    const first = deferred<LinkReference>();
    const latest = deferred<LinkReference>();
    list.mockResolvedValue(empty);
    create
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(latest.promise);
    renderView("a.invalid", 1, empty);
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));
    fireEvent.change(input(), { target: { value: "b.invalid" } });
    fireEvent.change(input(), { target: { value: "a.invalid" } });
    fireEvent.click(screen.getByRole("button", { name: "Add link" }));

    await act(async () => first.resolve(link));
    expect(screen.getByRole("button", { name: "Adding…" })).toBeDisabled();
    await act(async () => latest.resolve(link));
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled(),
    );
  });

  it("uses the active filter for pagination and resets to page one when it changes", async () => {
    list.mockImplementation(async (query: string, page: number) => ({
      ...all,
      current_page: page,
    }));
    renderView("design");
    fireEvent.click(results().getByRole("link", { name: "2" }));
    await waitFor(() => expect(list).toHaveBeenCalledWith("design", 2));

    fireEvent.change(input(), { target: { value: "browser" } });
    await waitFor(() => expect(list).toHaveBeenCalledWith("browser", 1));
    const update = onUrlUpdate.mock.lastCall![0];
    expect(update.searchParams.get("q")).toBe("browser");
    expect(update.searchParams.get("page")).toBeNull();
  });

  it("encodes query punctuation and completely clears an initial filtered page", async () => {
    renderView("original", 2);
    const query = "cats & dogs #pets";
    fireEvent.change(input(), { target: { value: query } });
    await waitFor(() => expect(list).toHaveBeenCalledWith(query, 1));
    expect(onUrlUpdate.mock.lastCall![0].searchParams.get("q")).toBe(query);

    fireEvent.click(screen.getByRole("button", { name: "Clear search" }));
    expect(input()).toHaveValue("");
    await waitFor(() => expect(list).toHaveBeenCalledWith("", 1));
    expect(onUrlUpdate.mock.lastCall![0].searchParams.toString()).toBe("");
  });

  it("never saves a URL loaded from the address bar", async () => {
    list.mockResolvedValue(single);
    renderView("a.invalid", 1, single);
    expect(input()).toHaveValue("a.invalid");
    expect(screen.getByRole("button", { name: "Saved" })).toBeDisabled();
    await waitFor(() => expect(list).toHaveBeenCalledWith("a.invalid", 1));
    expect(create).not.toHaveBeenCalled();
  });

  it("ignores a slower search response after the user changes the query", async () => {
    const slow = deferred<LinkListResult>();
    list.mockImplementation(async (query: string) => {
      if (query === "slow") {
        return slow.promise;
      }

      if (query === "fast") {
        return { ...single, links: [otherLink] };
      }

      return all;
    });
    renderView();
    fireEvent.change(input(), { target: { value: "slow" } });
    await waitFor(() => expect(list).toHaveBeenCalledWith("slow", 1));
    fireEvent.change(input(), { target: { value: "fast" } });
    await waitFor(() =>
      expect(results().getByRole("article", { name: "B link" })).toBeVisible(),
    );
    await act(async () => slow.resolve(single));

    expect(input()).toHaveValue("fast");
    expect(
      results().queryByRole("article", { name: "A link" }),
    ).not.toBeInTheDocument();
  });
});
