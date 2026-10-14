import { act, render, screen } from "@testing-library/react";
import type { ComponentType } from "react";
import { renderToString } from "react-dom/server";
import { expect, it, vi } from "vitest";

import { FeedPageBlocks } from "./FeedPageBlocks";

const { editor, loadEditor } = vi.hoisted(() => ({
  editor: {
    feed: { blocks: [{ type: "title" }, { type: "library", layout: "list" }] },
    isEditing: false,
    ready: undefined as Promise<void> | undefined,
  },
  loadEditor: vi.fn(),
}));

vi.mock("./Context", () => ({ useFeedBlockEditor: () => editor }));
vi.mock("./FeedBlockRender", () => ({
  FeedBlockRender: ({ block }: { block: { type: string } }) => (
    <article data-sd-block={block.type}>{block.type} content</article>
  ),
}));
vi.mock("./FeedPageEditor", () => ({
  FeedPageEditor: () => <p>Feed editing controls</p>,
}));
vi.mock("next/dynamic", async () => {
  const { lazy, Suspense } = await import("react");
  return {
    default: (
      load: () => Promise<ComponentType>,
      { loading: Loading }: { loading: ComponentType },
    ) => {
      const Editor = lazy(async () => {
        loadEditor();
        await editor.ready;
        return { default: await load() };
      });
      return () => (
        <Suspense fallback={<Loading />}>
          <Editor />
        </Suspense>
      );
    },
  };
});

it("renders the configured feed eagerly and retains it while editing controls load", async () => {
  const serverContent = document.createElement("div");
  serverContent.innerHTML = renderToString(<FeedPageBlocks />);
  expect(serverContent.textContent).toBe("title contentlibrary content");
  const view = render(<FeedPageBlocks />);
  expect(loadEditor).not.toHaveBeenCalled();
  expect(
    [...view.container.querySelectorAll("[data-sd-block]")].map((e) =>
      e.getAttribute("data-sd-block"),
    ),
  ).toEqual(["title", "library"]);

  let release!: () => void;
  editor.ready = new Promise((resolve) => {
    release = resolve;
  });
  editor.isEditing = true;
  view.rerender(<FeedPageBlocks />);
  expect(loadEditor).toHaveBeenCalledTimes(1);
  expect(screen.getByText("title content")).toBeVisible();
  expect(screen.getByText("library content")).toBeVisible();
  expect(screen.queryByText("Feed editing controls")).not.toBeInTheDocument();

  await act(async () => release());
  expect(await screen.findByText("Feed editing controls")).toBeVisible();
  editor.isEditing = false;
  view.rerender(<FeedPageBlocks />);
  expect(screen.getByText("title content")).toBeVisible();
  expect(screen.queryByText("Feed editing controls")).not.toBeInTheDocument();
});
