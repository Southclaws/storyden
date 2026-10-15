import { render, screen } from "@testing-library/react";
import { Suspense } from "react";
import { describe, expect, it, vi } from "vitest";

import { ContentEditor } from "./ContentEditor";

const editors = vi.hoisted(() => ({
  mode: undefined as "markdown" | "richtext" | undefined,
  rich: vi.fn(),
  markdown: vi.fn(),
}));

vi.mock("@/auth", () => ({
  useSession: () => ({ meta: { editor: { mode: editors.mode } } }),
}));

vi.mock("../ContentComposerRich/ContentComposerRich", () => {
  editors.rich();
  return { ContentComposerRich: () => <div>Rich editor</div> };
});

vi.mock("../ContentComposerMarkdown/ContentComposerMarkdown", () => {
  editors.markdown();
  return { ContentComposerMarkdown: () => <div>Markdown editor</div> };
});

describe("ContentEditor preference", () => {
  it("defaults to rich text and requests Markdown only when that preference is selected", async () => {
    const { rerender } = render(
      <Suspense fallback="Loading editor">
        <ContentEditor />
      </Suspense>,
    );
    expect(await screen.findByText("Rich editor")).toBeVisible();
    expect(editors.rich).toHaveBeenCalledOnce();
    expect(editors.markdown).not.toHaveBeenCalled();

    editors.mode = "markdown";
    rerender(
      <Suspense fallback="Loading editor">
        <ContentEditor />
      </Suspense>,
    );
    expect(await screen.findByText("Markdown editor")).toBeVisible();
    expect(editors.markdown).toHaveBeenCalledOnce();
    expect(screen.queryByText("Rich editor")).not.toBeInTheDocument();

    editors.mode = "richtext";
    rerender(
      <Suspense fallback="Loading editor">
        <ContentEditor />
      </Suspense>,
    );
    expect(await screen.findByText("Rich editor")).toBeVisible();
    expect(screen.queryByText("Markdown editor")).not.toBeInTheDocument();
  });
});
