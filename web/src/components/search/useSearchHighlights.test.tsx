import { act, cleanup, render, screen, waitFor } from "@testing-library/react";
import Link from "next/link";
import { type ReactNode, StrictMode } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { useSearchHighlights } from "./useSearchHighlights";

// jsdom supplies real DOM Ranges, but does not implement the painting API.
class TestHighlight extends Set<Range> {}

let registry: Map<string, TestHighlight>;

function Scope({
  query,
  enabled,
  children,
}: {
  query: string;
  enabled?: boolean;
  children: ReactNode;
}) {
  const ref = useSearchHighlights(query, enabled);
  return <div ref={ref}>{children}</div>;
}

function paintedText() {
  return [...(registry.get("storyden-search-match") ?? [])].map((range) =>
    range.toString(),
  );
}

beforeEach(() => {
  registry = new Map();
  vi.stubGlobal("Highlight", TestHighlight);
  vi.stubGlobal("CSS", { highlights: registry, supports: () => true });
});

afterEach(() => {
  cleanup();
  vi.unstubAllGlobals();
});

describe("useSearchHighlights", () => {
  it("paints repeated query words literally, without changing text or links", () => {
    const { container } = render(
      <Scope query="  garden GARDEN tools [x] c++  ">
        <Link data-search-highlight href="/garden">
          Garden tools, GARDEN [x] c++.
        </Link>
        <p>Garden metadata</p>
      </Scope>,
    );

    expect(paintedText()).toEqual(["Garden", "tools", "GARDEN", "[x]", "c++"]);
    expect(screen.getByRole("link")).toHaveAttribute("href", "/garden");
    expect(screen.getByRole("link").childNodes).toHaveLength(1);
    expect(container.querySelector("mark")).toBeNull();
  });

  it("uses original UTF-16 offsets after expanding case mappings and emoji", () => {
    render(
      <Scope query="garden café İstanbul">
        <p data-search-highlight>İstanbul 🐢 Garden CAFÉ garden</p>
      </Scope>,
    );
    expect(paintedText()).toEqual(["İstanbul", "Garden", "CAFÉ", "garden"]);
  });

  it("excludes hidden and editable text and does not join separate text nodes", () => {
    render(
      <Scope query="garden tools">
        <div data-search-highlight>
          Garden <span hidden>garden</span>
          <span aria-hidden="true">garden</span>
          <span contentEditable suppressContentEditableWarning>
            garden
          </span>
          <input defaultValue="garden" />
          <span>
            to<b>ols</b>
          </span>
          <span>tools</span>
        </div>
      </Scope>,
    );
    expect(paintedText()).toEqual(["Garden", "tools"]);
  });

  it("refreshes ranges when React replaces content without changing the query", async () => {
    const { rerender } = render(
      <Scope query="garden">
        <p data-search-highlight>Garden one</p>
      </Scope>,
    );
    const previous = [...registry.get("storyden-search-match")!];
    rerender(
      <Scope query="garden">
        <section data-search-highlight>New garden garden</section>
      </Scope>,
    );
    await waitFor(() => expect(paintedText()).toEqual(["garden", "garden"]));
    expect([...registry.get("storyden-search-match")!]).not.toContain(
      previous[0],
    );
  });

  it("removes only its own ranges when another scope updates or unmounts", () => {
    const first = render(
      <Scope query="garden">
        <p data-search-highlight>Garden</p>
      </Scope>,
    );
    const second = render(
      <Scope query="tools">
        <p data-search-highlight>Tools shed</p>
      </Scope>,
    );
    const unrelated = new TestHighlight();
    registry.set("other-feature", unrelated);
    first.rerender(
      <Scope query="shed">
        <p data-search-highlight>Garden shed</p>
      </Scope>,
    );
    expect(paintedText()).toEqual(["Tools", "shed"]);
    first.unmount();
    expect(paintedText()).toEqual(["Tools"]);
    second.unmount();
    expect(registry.has("storyden-search-match")).toBe(false);
    expect(registry.get("other-feature")).toBe(unrelated);
  });

  it("cancels queued content updates when the scope unmounts", async () => {
    const { rerender, unmount } = render(
      <Scope query="garden">
        <p data-search-highlight>Garden</p>
      </Scope>,
    );
    const cancelFrame = vi.spyOn(window, "cancelAnimationFrame");
    await act(async () => {
      rerender(
        <Scope query="garden">
          <p data-search-highlight>Garden garden</p>
        </Scope>,
      );
    });
    unmount();
    expect(cancelFrame).toHaveBeenCalledOnce();
    expect(registry.has("storyden-search-match")).toBe(false);
  });

  it("cleans up for empty queries, editing, and Strict Mode remounts", () => {
    const view = (query: string, enabled = true) => (
      <StrictMode>
        <Scope query={query} enabled={enabled}>
          <p data-search-highlight>Garden</p>
        </Scope>
      </StrictMode>
    );
    const { rerender, unmount } = render(view("garden"));
    expect(paintedText()).toEqual(["Garden"]);
    rerender(view("  "));
    expect(paintedText()).toEqual([]);
    rerender(view("garden", false));
    expect(paintedText()).toEqual([]);
    rerender(view("garden"));
    expect(paintedText()).toEqual(["Garden"]);
    unmount();
    expect(registry.size).toBe(0);
  });

  it.each(["Highlight", "CSS", "registry", "selector"])(
    "preserves content when %s is unsupported",
    (missing) => {
      if (missing === "Highlight" || missing === "CSS") {
        vi.stubGlobal(missing, undefined);
      } else {
        vi.stubGlobal("CSS", {
          highlights: missing === "registry" ? undefined : registry,
          supports: () => missing !== "selector",
        });
      }
      render(
        <Scope query="garden">
          <p data-search-highlight>Garden</p>
        </Scope>,
      );
      expect(screen.getByText("Garden")).toBeVisible();
      expect(registry.size).toBe(0);
    },
  );
});
