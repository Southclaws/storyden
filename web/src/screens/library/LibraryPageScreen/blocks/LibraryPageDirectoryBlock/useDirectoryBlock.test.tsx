import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { getDefaultBlockConfig } from "./column";
import { useDirectoryBlock } from "./useDirectoryBlock";

const overwriteBlock = vi.fn();

let blockMock: any;
let schemaMock: any[];

vi.mock("../useBlock", () => ({
  useBlock: () => blockMock,
}));

vi.mock("../../store", () => ({
  useWatch: (selector: (state: any) => unknown) =>
    selector({
      draft: {
        child_property_schema: schemaMock,
      },
    }),
}));

vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: {
      getState: () => ({
        overwriteBlock,
      }),
    },
  }),
}));

describe("useDirectoryBlock", () => {
  beforeEach(() => {
    overwriteBlock.mockReset();
    schemaMock = [
      {
        fid: "p1",
        name: "Priority",
        sort: "0",
        type: "text",
      },
    ];
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("returns and writes a default config when missing", async () => {
    blockMock = {
      type: "directory",
      config: undefined,
    };

    const { result } = renderHook(() => useDirectoryBlock());
    const expected = getDefaultBlockConfig(schemaMock as any);

    expect(result.current.config).toEqual(expected);

    await act(async () => {
      vi.runOnlyPendingTimers();
    });

    expect(overwriteBlock).toHaveBeenCalledWith({
      type: "directory",
      config: expected,
    });
  });

  it("does not rewrite config when schema and config already match", async () => {
    blockMock = {
      type: "directory",
      config: getDefaultBlockConfig(schemaMock as any),
    };

    renderHook(() => useDirectoryBlock());

    await act(async () => {
      vi.runOnlyPendingTimers();
    });

    expect(overwriteBlock).not.toHaveBeenCalled();
  });
});
