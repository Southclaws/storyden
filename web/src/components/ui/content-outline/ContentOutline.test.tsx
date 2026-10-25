import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderToString } from "react-dom/server";
import {
  afterEach,
  beforeAll,
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from "vitest";

import { ContentComposerMarkdown } from "@/components/content/ContentComposerMarkdown/ContentComposerMarkdown";
import { ContentComposerRich } from "@/components/content/ContentComposerRich/ContentComposerRich";
import { collectLibraryOutline } from "@/screens/library/LibraryPageScreen/outline";
import { collectThreadOutline } from "@/screens/thread/ThreadScreen/outline";

import { ContentOutline } from "./ContentOutline";
import {
  collectContentHeadings,
  getActiveOutlineKey,
  scrollToOutlineEntry,
} from "./outline";

let desktop = true;
let reducedMotion = false;
let mediaChange: (() => void) | undefined;
let resizeCallbacks: Map<Element, ResizeObserverCallback>;

beforeAll(async () => {
  // Keep the real rail, but load its dependencies outside assertion timeouts.
  await import("./ContentOutline.internal");
});

beforeEach(() => {
  vi.stubGlobal(
    "IntersectionObserver",
    class {
      observe() {}
      disconnect() {}
      unobserve() {}
    },
  );
  desktop = true;
  reducedMotion = false;
  mediaChange = undefined;
  resizeCallbacks = new Map();
  vi.stubGlobal("innerHeight", 900);
  vi.stubGlobal(
    "ResizeObserver",
    class {
      constructor(private callback: ResizeObserverCallback) {}
      observe(target: Element) {
        resizeCallbacks.set(target, this.callback);
      }
      disconnect() {}
      unobserve() {}
    },
  );
  vi.stubGlobal(
    "matchMedia",
    vi.fn((query: string) => ({
      get matches() {
        return query.includes("reduced-motion") ? reducedMotion : desktop;
      },
      media: query,
      addEventListener: (_: string, callback: () => void) => {
        mediaChange = callback;
      },
      removeEventListener: vi.fn(),
    })),
  );
  vi.spyOn(HTMLElement.prototype, "getClientRects").mockReturnValue([
    {},
  ] as unknown as DOMRectList);
  vi.spyOn(HTMLElement.prototype, "getBoundingClientRect").mockImplementation(
    function (this: HTMLElement) {
      const top = Number(this.dataset["top"] ?? 100) - window.scrollY;
      const height = Number(
        this.dataset["height"] ??
          (this.classList.contains("content-outline__body") ? 2000 : 20),
      );
      return {
        top,
        height,
        bottom: top + height,
        left: 0,
        right: 100,
        width: 100,
        x: 0,
        y: top,
        toJSON() {},
      };
    },
  );
  vi.spyOn(window, "scrollTo").mockImplementation(() => {});
});

afterEach(() => vi.unstubAllGlobals());

function collect(root: HTMLElement) {
  return collectContentHeadings(root, "story");
}

function Example({
  disabled = false,
  extra = false,
}: {
  disabled?: boolean;
  extra?: boolean;
}) {
  return (
    <ContentOutline
      label="On this page"
      collectEntries={collect}
      disabled={disabled}
    >
      <h2 data-top="0">Introduction</h2>
      <p>Some content.</p>
      <h3 data-top="900">Details</h3>
      {extra && <h2>Conclusion</h2>}
      <button>Outside</button>
    </ContentOutline>
  );
}

describe("ContentOutline", () => {
  it.each([
    [1799, false],
    [1800, true],
    [1801, true],
  ])(
    "requires two viewport heights (%ipx, visible=%s)",
    async (height, visible) => {
      const view = render(<Example />);
      const body = view.container.querySelector<HTMLElement>(
        ".content-outline__body",
      )!;
      body.dataset["height"] = String(height);
      expect(view.container.firstChild).toHaveAttribute(
        "data-visible",
        "false",
      );
      await act(() => new Promise(requestAnimationFrame));
      expect(view.container.firstChild).toHaveAttribute(
        "data-visible",
        String(visible),
      );
    },
  );

  it("keeps reading content in server HTML without rendering a rail", () => {
    const html = renderToString(<Example />);
    expect(html).toContain("Introduction");
    expect(html).toContain("Some content.");
    expect(html).toContain('data-visible="false"');
    expect(html).not.toContain("data-content-outline");
  });

  it("remeasures unchanged entries on viewport, content and late layout changes", async () => {
    const view = render(<Example />);
    await screen.findByRole("button", { name: "On this page" });
    const body = view.container.querySelector<HTMLElement>(
      ".content-outline__body",
    )!;

    vi.stubGlobal("innerHeight", 1200);
    fireEvent.resize(window);
    await waitFor(() =>
      expect(view.container.firstChild).toHaveAttribute(
        "data-visible",
        "false",
      ),
    );

    vi.stubGlobal("innerHeight", 900);
    fireEvent.resize(window);
    await screen.findByRole("button", { name: "On this page" });

    body.dataset["height"] = "1500";
    act(() => resizeCallbacks.get(body)!([], {} as ResizeObserver));
    await waitFor(() =>
      expect(view.container.firstChild).toHaveAttribute(
        "data-visible",
        "false",
      ),
    );

    body.dataset["height"] = "2000";
    body.querySelector("p")!.textContent = "More reading content.";
    await screen.findByRole("button", { name: "On this page" });

    body.dataset["height"] = "1200";
    body.querySelector("p")!.textContent = "Shortened content.";
    await waitFor(() =>
      expect(view.container.firstChild).toHaveAttribute(
        "data-visible",
        "false",
      ),
    );
  });

  it("counts original content, replies and spacing, excluding the reply composer", async () => {
    const view = render(
      <ContentOutline
        label="In this thread"
        collectEntries={collectThreadOutline}
      >
        <div data-outline-original="" data-top="100" data-height="800">
          Original content
        </div>
        <div
          data-outline-reply=""
          id="reply-1"
          data-top="1000"
          data-height="500"
        />
        <div data-top="1500" data-height="2000">
          Reply composer
        </div>
      </ContentOutline>,
    );
    await act(() => new Promise(requestAnimationFrame));
    expect(view.container.firstChild).toHaveAttribute("data-visible", "false");

    const reply = view.container.querySelector<HTMLElement>(
      "[data-outline-reply]",
    )!;
    reply.dataset["height"] = "900";
    reply.textContent = "A longer reply.";
    await screen.findByRole("button", { name: "In this thread" });
    expect(view.container.firstChild).toHaveAttribute("data-visible", "true");

    vi.spyOn(window, "scrollY", "get").mockReturnValue(300);
    fireEvent.scroll(window);
    fireEvent.resize(window);
    await act(() => new Promise(requestAnimationFrame));
    expect(view.container.firstChild).toHaveAttribute("data-visible", "true");
  });

  it("measures Library reading blocks without counting the cover or title", async () => {
    const view = render(
      <ContentOutline
        label="On this page"
        collectEntries={collectLibraryOutline}
      >
        <div data-outline-block="cover" data-top="0" data-height="2000">
          Cover
        </div>
        <div data-outline-block="title" data-top="2000" data-height="100">
          Title
        </div>
        <div data-outline-block="content" data-top="2100" data-height="800">
          <h2>First section</h2>
          <h3>Details</h3>
        </div>
        <div data-outline-block="directory" data-top="3000" data-height="300">
          Directory
        </div>
      </ContentOutline>,
    );
    await act(() => new Promise(requestAnimationFrame));
    expect(view.container.firstChild).toHaveAttribute("data-visible", "false");

    const directory = view.container.querySelector<HTMLElement>(
      '[data-outline-block="directory"]',
    )!;
    directory.dataset["height"] = "900";
    directory.textContent = "More pages.";
    await screen.findByRole("button", { name: "On this page" });
    expect(view.container.firstChild).toHaveAttribute("data-visible", "true");
  });

  it("marks only overflowing ends, including after scrolling and content resizing", async () => {
    render(<Example />);
    const trigger = await screen.findByRole("button", { name: "On this page" });
    const ticks = trigger.querySelector<HTMLElement>('[aria-hidden="true"]')!;
    Object.defineProperties(ticks, {
      scrollHeight: { configurable: true, value: 500 },
      clientHeight: { configurable: true, value: 300 },
    });
    fireEvent.scroll(ticks);
    await waitFor(() => {
      expect(ticks).toHaveAttribute("data-overflow-block-start", "false");
      expect(ticks).toHaveAttribute("data-overflow-block-end", "true");
    });
    ticks.scrollTop = 100;
    fireEvent.scroll(ticks);
    await waitFor(() => {
      expect(ticks).toHaveAttribute("data-overflow-block-start", "true");
      expect(ticks).toHaveAttribute("data-overflow-block-end", "true");
    });
    ticks.scrollTop = 200;
    fireEvent.scroll(ticks);
    await waitFor(() => {
      expect(ticks).toHaveAttribute("data-overflow-block-start", "true");
      expect(ticks).toHaveAttribute("data-overflow-block-end", "false");
    });
    Object.defineProperty(ticks, "scrollHeight", { value: 300 });
    ticks.scrollTop = 0;
    // A content mutation triggers the same measurement as ResizeObserver.
    ticks.appendChild(document.createElement("span"));
    await waitFor(() => {
      expect(ticks).toHaveAttribute("data-overflow-block-start", "false");
      expect(ticks).toHaveAttribute("data-overflow-block-end", "false");
    });
  });

  it("shows reply author, relative date and snippet and updates reply metadata", async () => {
    const createdAt = new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString();
    const view = render(
      <ContentOutline
        label="In this thread"
        collectEntries={collectThreadOutline}
      >
        <div data-outline-original="" data-height="1800">
          Original content
        </div>
        <div
          data-outline-reply=""
          id="reply-1"
          data-outline-label="Ada"
          data-outline-description="A useful observation"
          data-outline-created-at={createdAt}
        />
      </ContentOutline>,
    );
    fireEvent.pointerEnter(
      await screen.findByRole("button", { name: "In this thread" }),
    );
    const row = await screen.findByRole("button", {
      name: /Ada.*2 hours ago.*A useful observation/,
    });
    expect(row.querySelector("time")).toHaveAttribute("datetime", createdAt);
    expect(screen.getByText("Replies", { exact: true })).toBeVisible();
    expect(row.querySelector("time")!.parentElement!.textContent).toContain(
      "Ada",
    );
    const updatedAt = new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString();
    view.container
      .querySelector("[data-outline-reply]")!
      .setAttribute("data-outline-created-at", updatedAt);
    await screen.findByRole("button", {
      name: /Ada.*4 hours ago.*A useful observation/,
    });
  });

  it.each([ContentComposerMarkdown, ContentComposerRich])(
    "navigates rendered composer headings without changing saved content (%#)",
    async (Composer) => {
      const onChange = vi.fn();
      const user = userEvent.setup({ skipHover: true });
      render(
        <ContentOutline label="On this page" collectEntries={collect}>
          <Composer
            disabled
            onChange={onChange}
            initialValue="<h2>Overview</h2><p>First paragraph.</p><h3>Details with <strong>emphasis</strong></h3><p>Second paragraph.</p>"
          />
        </ContentOutline>,
      );
      await user.click(
        await screen.findByRole("button", { name: "On this page" }),
      );
      await user.click(
        await screen.findByRole("button", { name: "Details with emphasis" }),
      );
      expect(
        screen.getByRole("heading", { name: "Details with emphasis" }),
      ).toHaveFocus();
      expect(onChange).not.toHaveBeenCalled();
    },
  );
  it("opens on hover, lets the pointer cross to the panel, and dismisses on leaving", async () => {
    render(<Example />);
    const trigger = await screen.findByRole("button", { name: "On this page" });
    fireEvent.pointerEnter(trigger);
    const panel = await screen.findByRole("dialog");
    expect(document.activeElement).not.toBe(trigger);
    fireEvent.pointerLeave(trigger);
    fireEvent.pointerEnter(panel);
    expect(panel).toBeVisible();
    fireEvent.pointerLeave(panel);
    await waitFor(() =>
      expect(trigger).toHaveAttribute("aria-expanded", "false"),
    );
  });

  it("keeps the first click open after focus, then closes and reopens on repeated clicks", async () => {
    const user = userEvent.setup({ skipHover: true });
    render(<Example />);

    const trigger = await screen.findByRole("button", { name: "On this page" });

    await user.click(trigger);
    expect(trigger).toHaveFocus();
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByRole("dialog")).toBeVisible();

    await user.click(trigger);
    expect(trigger).toHaveFocus();
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    await waitFor(() => expect(screen.queryByRole("dialog")).toBeNull());

    await user.click(trigger);
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByRole("dialog")).toBeVisible();
  });

  it("closes a hover-opened outline when its trigger is clicked", async () => {
    const user = userEvent.setup();
    render(<Example />);

    const trigger = await screen.findByRole("button", { name: "On this page" });

    await user.hover(trigger);
    expect(trigger).toHaveAttribute("aria-expanded", "true");

    await user.click(trigger);
    expect(trigger).toHaveFocus();
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    await waitFor(() => expect(screen.queryByRole("dialog")).toBeNull());
  });

  it("opens, closes and reopens on repeated touch activation", async () => {
    const user = userEvent.setup();
    render(<Example />);

    const trigger = await screen.findByRole("button", { name: "On this page" });

    await user.pointer({ keys: "[TouchA>]", target: trigger });
    await user.pointer({ keys: "[/TouchA]", target: trigger });
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByRole("dialog")).toBeVisible();

    await user.pointer({ keys: "[TouchA>]", target: trigger });
    await user.pointer({ keys: "[/TouchA]", target: trigger });
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    await waitFor(() => expect(screen.queryByRole("dialog")).toBeNull());

    await user.pointer({ keys: "[TouchA>]", target: trigger });
    await user.pointer({ keys: "[/TouchA]", target: trigger });
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByRole("dialog")).toBeVisible();
  });

  it("supports keyboard navigation, Escape, re-entry and destination focus", async () => {
    const user = userEvent.setup();
    render(<Example />);
    const trigger = await screen.findByRole("button", { name: "On this page" });
    trigger.focus();
    await user.keyboard("{Enter}");
    const first = await screen.findByRole("button", { name: "Introduction" });
    await waitFor(() => expect(first).toHaveFocus());
    await user.keyboard("{Tab}");
    expect(screen.getByRole("button", { name: "Details" })).toHaveFocus();
    await user.keyboard("{Escape}");
    await waitFor(() =>
      expect(trigger).toHaveAttribute("aria-expanded", "false"),
    );
    await waitFor(() => expect(trigger).toHaveFocus());
    await user.keyboard("{Enter}");
    await user.click(await screen.findByRole("button", { name: "Details" }));
    await waitFor(() =>
      expect(screen.getByRole("heading", { name: "Details" })).toHaveFocus(),
    );
    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 828,
      behavior: "smooth",
    });
    expect(trigger).toHaveAttribute("aria-expanded", "false");
  });

  it("updates after content changes and suppresses the rail while editing or on mobile", async () => {
    const view = render(<Example />);
    await screen.findByRole("button", { name: "On this page" });
    view.rerender(<Example extra />);
    fireEvent.pointerEnter(
      screen.getByRole("button", { name: "On this page" }),
    );
    expect(
      await screen.findByRole("button", { name: "Conclusion" }),
    ).toBeVisible();
    view.rerender(<Example disabled />);
    expect(screen.queryByRole("button", { name: "On this page" })).toBeNull();
    view.rerender(<Example />);
    await screen.findByRole("button", { name: "On this page" });
    desktop = false;
    mediaChange?.();
    await waitFor(() =>
      expect(screen.queryByRole("button", { name: "On this page" })).toBeNull(),
    );
  });

  it("does not render for a single heading", async () => {
    render(
      <ContentOutline label="On this page" collectEntries={collect}>
        <h2>Only section</h2>
      </ContentOutline>,
    );
    await new Promise((resolve) => requestAnimationFrame(resolve));
    expect(screen.queryByRole("button", { name: "On this page" })).toBeNull();
  });

  it("dismisses on outside interaction without returning focus to the trigger", async () => {
    const user = userEvent.setup({ skipHover: true });
    render(<Example />);
    const trigger = await screen.findByRole("button", { name: "On this page" });
    await user.click(trigger);
    expect(await screen.findByRole("dialog")).toBeVisible();
    await user.click(screen.getByRole("button", { name: "Outside" }));
    await waitFor(() =>
      expect(trigger).toHaveAttribute("aria-expanded", "false"),
    );
    expect(screen.getByRole("button", { name: "Outside" })).toHaveFocus();
  });
});

describe("outline sources", () => {
  it("keeps reading blocks in displayed order, ignores empty/decorative blocks and scopes headings to content", () => {
    const { container } = render(
      <>
        <div data-outline-block="cover">Cover</div>
        <div data-outline-block="title">
          <h1>Page title</h1>
        </div>
        <div data-outline-block="directory">
          <h2>A card title</h2>
        </div>
        <div data-outline-block="properties" data-height="0" />
        <div data-outline-block="content">
          <h2>Start</h2>
          <h4>Repeated</h4>
          <h3>Repeated</h3>
          <h2>End</h2>
          <h3> </h3>
          <div data-no-typography>
            <h2>Embedded preview title</h2>
          </div>
        </div>
        <div data-outline-block="tags">Tag</div>
        <div data-outline-block="assets">
          <img src="/test-asset.png" alt="Example" />
        </div>
      </>,
    );
    const entries = collectLibraryOutline(container);
    expect(entries.map(({ label, depth }) => [label, depth])).toEqual([
      ["Directory", 0],
      ["Content", 0],
      ["Start", 1],
      ["Repeated", 2],
      ["Repeated", 2],
      ["End", 1],
      ["Assets", 0],
    ]);
    expect(new Set(entries.map((e) => e.key)).size).toBe(7);
    expect(entries[3]!.target).not.toBe(entries[4]!.target);
  });

  it("includes first-post headings only above one viewport and never includes reply headings", () => {
    const view = render(
      <>
        <div data-outline-original="">
          <h1>Thread</h1>
          <div data-outline-body="" data-height="2000">
            <h2>Long section</h2>
          </div>
        </div>
        <div
          data-outline-reply=""
          id="reply-51"
          data-outline-label="Ada"
          data-outline-description="A useful observation"
        >
          <h2>Ignore reply heading</h2>
        </div>
      </>,
    );
    expect(collectThreadOutline(view.container).map((e) => e.label)).toEqual([
      "Original post",
      "Long section",
      "Ada",
    ]);
    view.container.querySelector<HTMLElement>("[data-outline-body]")!.dataset[
      "height"
    ] = String(window.innerHeight);
    expect(collectThreadOutline(view.container).map((e) => e.key)).toEqual([
      "original",
      "reply-51",
    ]);
  });

  it("tracks the last crossed target including the document bottom", () => {
    const { container } = render(
      <>
        <h2 data-top="-200">First</h2>
        <h2 data-top="80.25">Middle</h2>
        <h2 data-top="900">Last</h2>
      </>,
    );
    const entries = collect(container);
    expect(getActiveOutlineKey(entries, 80)).toBe("story:heading:1");
    vi.spyOn(window, "scrollY", "get").mockReturnValue(400);
    vi.spyOn(document.documentElement, "scrollHeight", "get").mockReturnValue(
      400 + window.innerHeight,
    );
    expect(getActiveOutlineKey(entries, 80)).toBe("story:heading:2");
  });

  it("respects reduced motion and ignores detached destinations", () => {
    const { container, unmount } = render(<h2 data-top="700">Target</h2>);
    const entry = collect(container)[0]!;
    reducedMotion = true;
    scrollToOutlineEntry(entry);
    expect(window.scrollTo).toHaveBeenCalledWith({
      top: 628,
      behavior: "instant",
    });
    unmount();
    scrollToOutlineEntry(entry);
    expect(window.scrollTo).toHaveBeenCalledTimes(1);
  });
});
