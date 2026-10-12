import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { ComponentType, PropsWithChildren } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { CommandPalette } from "./CommandPalette";

const { auth, loadContent } = vi.hoisted(() => ({
  auth: { authenticated: false },
  loadContent: vi.fn(),
}));

vi.mock("@/auth", () => ({
  useSession: () => (auth.authenticated ? {} : undefined),
}));

vi.mock("@/lib/settings/capabilities", () => ({
  useCapability: () => true,
}));

vi.mock("@/utils/permissions", () => ({
  hasPermission: () => true,
}));

vi.mock("@radix-ui/react-dialog", () => ({
  Title: ({ children }: PropsWithChildren) => children,
}));

vi.mock("cmdk", () => ({
  Command: {
    Dialog: ({ open, children }: PropsWithChildren<{ open: boolean }>) =>
      open ? <div role="dialog">{children}</div> : null,
  },
}));

vi.mock("next/dynamic", async () => {
  const { lazy, Suspense } = await import("react");

  return {
    default: (load: () => Promise<ComponentType>, { loading: Loading }) => {
      const Content = lazy(async () => {
        loadContent();
        return { default: await load() };
      });

      return () => (
        <Suspense fallback={<Loading />}>
          <Content />
        </Suspense>
      );
    },
  };
});

vi.mock("./CommandPaletteDialogContent", async () => {
  const { useCommandPalette } = await import("./Context");

  return {
    CommandPaletteDialogContent: () => {
      const { initialSessionID, loadChatSession, search, setSearch } =
        useCommandPalette();

      return (
        <>
          <label>
            Search
            <input
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>
          <button onClick={() => loadChatSession("session-123")}>
            Open saved chat
          </button>
          {initialSessionID && <p>Selected chat: {initialSessionID}</p>}
        </>
      );
    },
  };
});

describe("CommandPalette loading boundary", () => {
  beforeEach(() => {
    auth.authenticated = false;
    loadContent.mockClear();
  });

  it("waits for authenticated activation and retains selected chat across reopening", async () => {
    const view = render(<CommandPalette />);

    expect(loadContent).not.toHaveBeenCalled();
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

    fireEvent.keyDown(document, { key: "k", metaKey: true });
    expect(loadContent).not.toHaveBeenCalled();
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

    auth.authenticated = true;
    view.rerender(<CommandPalette />);
    fireEvent.keyDown(document, { key: "k", metaKey: true });

    expect(screen.getByRole("status", { name: "Loading" })).toBeVisible();
    const search = await screen.findByRole("textbox", { name: "Search" });
    expect(loadContent).toHaveBeenCalledTimes(1);

    fireEvent.change(search, { target: { value: "find a page" } });
    fireEvent.keyDown(document, { key: "k", metaKey: true });
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

    fireEvent.keyDown(document, { key: "k", ctrlKey: true });
    expect(await screen.findByRole("textbox", { name: "Search" })).toHaveValue(
      "find a page",
    );
    fireEvent.click(screen.getByRole("button", { name: "Open saved chat" }));
    expect(screen.getByText("Selected chat: session-123")).toBeVisible();

    // Escape first leaves chat mode, then closes the menu.
    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.getByRole("dialog")).toBeVisible();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();

    fireEvent.keyDown(document, { key: "k", ctrlKey: true });
    await waitFor(() => {
      expect(screen.getByText("Selected chat: session-123")).toBeVisible();
    });
    expect(loadContent).toHaveBeenCalledTimes(1);
  });
});
