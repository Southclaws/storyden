import { act, render, screen } from "@testing-library/react";
import { renderToString } from "react-dom/server";
import { expect, it, vi } from "vitest";

import type { Account } from "@/api/openapi-schema";

import { NavigationItems, type Props } from "./NavigationItems";

const { loadEditor, loading } = vi.hoisted(() => ({
  loadEditor: vi.fn(),
  loading: { ready: undefined as Promise<void> | undefined },
}));

vi.mock("@/components/category/CategoryList/CategoryList", () => ({
  CategoryList: () => null,
}));
vi.mock("../LibraryNavigationTree/LibraryNavigationTree", () => ({
  LibraryNavigationTree: () => null,
}));
vi.mock("./NavigationItemsEditor", async () => {
  loadEditor();
  await loading.ready;
  return {
    NavigationItemsEditor: () => <button>Edit navigation</button>,
  };
});

it("keeps server-rendered navigation available until the editing controls are ready", async () => {
  const props: Props = {
    navigation: {
      items: [
        { type: "members" },
        { type: "custom-link", id: "guide", label: "Guide", href: "/guide" },
        { type: "robots" },
      ],
    },
    currentPath: "/",
    isEditing: false,
    robotsEnabled: false,
  };
  const serverContent = document.createElement("div");
  serverContent.innerHTML = renderToString(<NavigationItems {...props} />);
  expect(
    [...serverContent.querySelectorAll("a")].map((a) => a.getAttribute("href")),
  ).toEqual(["/m", "/guide"]);

  const view = render(<NavigationItems {...props} />);
  expect(loadEditor).not.toHaveBeenCalled();
  expect(screen.getByRole("link", { name: "Guide" })).toHaveAttribute(
    "href",
    "/guide",
  );
  expect(
    screen.queryByRole("link", { name: "Robots" }),
  ).not.toBeInTheDocument();

  let release!: () => void;
  loading.ready = new Promise((resolve) => {
    release = resolve;
  });
  view.rerender(<NavigationItems {...props} isEditing />);
  expect(screen.getByRole("link", { name: "Guide" })).toBeVisible();
  expect(
    screen.queryByRole("button", { name: "Edit navigation" }),
  ).not.toBeInTheDocument();
  await act(async () => release());
  expect(
    await screen.findByRole("button", { name: "Edit navigation" }),
  ).toBeVisible();
  expect(loadEditor).toHaveBeenCalledTimes(1);

  view.rerender(<NavigationItems {...props} robotsEnabled />);
  expect(screen.getByRole("link", { name: "Guide" })).toBeVisible();
  expect(screen.getByRole("link", { name: "Robots" })).toBeVisible();
  expect(
    screen.queryByRole("button", { name: "Edit navigation" }),
  ).not.toBeInTheDocument();
});

it("renders Unread only for signed-in members", () => {
  const props: Props = {
    navigation: { items: [{ type: "unread" }] },
    currentPath: "/unread",
    isEditing: false,
    robotsEnabled: false,
  };

  const view = render(<NavigationItems {...props} />);
  expect(screen.queryByRole("link", { name: "Unread" })).toBeNull();

  view.rerender(
    <NavigationItems {...props} session={{ id: "member-id" } as Account} />,
  );
  expect(screen.getByRole("link", { name: "Unread" })).toHaveAttribute(
    "href",
    "/unread",
  );
});
