import { render } from "@testing-library/react";
import { PropsWithChildren } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { FeedConfig } from "@/lib/settings/feed";

import { InitialData } from "../types";

import {
  FeedBlockEditorActions,
  FeedBlockEditorProvider,
  useFeedBlockEditor,
} from "./Context";

const mocks = vi.hoisted(() => ({
  updateFeed: vi.fn(),
  updateSettings: vi.fn(),
}));

vi.mock("@/lib/settings/feed-client", () => ({
  useFeedMutation: () => ({ updateFeed: mocks.updateFeed }),
}));

vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ updateSettings: mocks.updateSettings }),
}));

const INITIAL_FEED = {
  blocks: [{ type: "title" }, { type: "content" }],
} satisfies FeedConfig;

describe("FeedBlockEditorProvider mutations", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("serializes direct editor mutations against the saved evolving feed", async () => {
    let releaseFirstSave = () => {};
    const firstSave = new Promise<void>((resolve) => {
      releaseFirstSave = resolve;
    });
    mocks.updateFeed
      .mockImplementationOnce(() => firstSave)
      .mockResolvedValueOnce(undefined);
    const editor = renderEditor();

    const add = editor.addBlock("cover");
    const remove = editor.removeBlock("title");

    await vi.waitFor(() => expect(mocks.updateFeed).toHaveBeenCalledOnce());
    expect(mocks.updateFeed).toHaveBeenNthCalledWith(1, {
      blocks: [{ type: "title" }, { type: "content" }, { type: "cover" }],
    });

    releaseFirstSave();
    const [added, removed] = await Promise.all([add, remove]);

    expect(added.blocks.map((block) => block.type)).toEqual([
      "title",
      "content",
      "cover",
    ]);
    expect(removed.blocks.map((block) => block.type)).toEqual([
      "content",
      "cover",
    ]);
    expect(mocks.updateFeed).toHaveBeenNthCalledWith(2, {
      blocks: [{ type: "content" }, { type: "cover" }],
    });
  });

  it("rolls a failed mutation back before starting the next queued mutation", async () => {
    let rejectFirstSave = (_error: Error) => {};
    const firstSave = new Promise<void>((_resolve, reject) => {
      rejectFirstSave = reject;
    });
    mocks.updateFeed
      .mockImplementationOnce(() => firstSave)
      .mockResolvedValueOnce(undefined);
    const editor = renderEditor();

    const add = editor.addBlock("cover");
    const remove = editor.removeBlock("title");

    await vi.waitFor(() => expect(mocks.updateFeed).toHaveBeenCalledOnce());
    rejectFirstSave(new Error("save failed"));

    await expect(add).rejects.toThrow("save failed");
    await expect(remove).resolves.toEqual({ blocks: [{ type: "content" }] });
    expect(mocks.updateFeed).toHaveBeenNthCalledWith(2, {
      blocks: [{ type: "content" }],
    });
  });
});

function EditorProbe({
  onEditor,
}: {
  onEditor: (editor: FeedBlockEditorActions) => void;
}) {
  onEditor(useFeedBlockEditor());
  return null;
}

function Provider({ children }: PropsWithChildren) {
  return (
    <FeedBlockEditorProvider
      feed={structuredClone(INITIAL_FEED)}
      initialData={{} as InitialData}
      isEditing
    >
      {children}
    </FeedBlockEditorProvider>
  );
}

function renderEditor() {
  let editor: FeedBlockEditorActions | undefined;
  render(
    <Provider>
      <EditorProbe onEditor={(value) => (editor = value)} />
    </Provider>,
  );
  expect(editor).toBeDefined();
  return editor!;
}
