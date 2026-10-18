import { act, renderHook } from "@testing-library/react";
import type { FixedCropperRef } from "react-advanced-cropper";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { assetUpload } from "@/api/openapi-client/assets";
import type { Asset } from "@/api/openapi-schema";

import { useLibraryPageCoverBlock } from "./useLibraryPageCoverBlock";

let primaryImage: Asset | undefined;
const setPrimaryImage = vi.fn();

vi.mock("@/api/openapi-client/assets", () => ({ assetUpload: vi.fn() }));
vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<unknown>) => fn(),
}));
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: {
      getState: () => ({
        draft: { primary_image: primaryImage },
        setPrimaryImage,
      }),
    },
  }),
}));

const original: Asset = {
  id: "original",
  filename: "original.png",
  path: "/api/assets/original.png",
  mime_type: "image/png",
  width: 1800,
  height: 900,
};
const crop: Asset = {
  id: "crop",
  filename: "crop.png",
  path: "/api/assets/crop.png",
  mime_type: "image/png",
  width: 1536,
  height: 384,
  parent: original,
};

beforeEach(() => {
  primaryImage = undefined;
  setPrimaryImage.mockReset();
});

describe("Library cover cropping", () => {
  it.each(["empty", "existing"])(
    "uses a newly uploaded original after mounting with an %s cover",
    async (initial) => {
      primaryImage = initial === "existing" ? crop : undefined;
      const { result } = renderHook(() => useLibraryPageCoverBlock());
      const replacement = { ...original, id: "replacement" };
      primaryImage = replacement;

      const blob = new Blob(["cropped pixels"], { type: "image/png" });
      const canvas = document.createElement("canvas");
      vi.spyOn(canvas, "toBlob").mockImplementation((callback) =>
        callback(blob),
      );
      const coordinates = { left: 0, top: 100, width: 1800, height: 450 };
      result.current.cropperRef.current = {
        getCanvas: () => canvas,
        getCoordinates: () => coordinates,
      } as FixedCropperRef;
      vi.mocked(assetUpload).mockResolvedValue({
        ...crop,
        parent: replacement,
      });

      await act(async () => result.current.handleInteractionEnd());

      expect(assetUpload).toHaveBeenCalledWith(blob, {
        filename: "cropped-cover",
        parent_asset_id: "replacement",
      });
      expect(setPrimaryImage).toHaveBeenCalledWith({
        isReplacement: false,
        asset: { ...crop, parent: replacement },
        config: coordinates,
      });
    },
  );

  it("keeps repeated crops as siblings of the original", async () => {
    primaryImage = crop;
    const { result } = renderHook(() => useLibraryPageCoverBlock());
    const blob = new Blob(["second crop"], { type: "image/png" });
    const canvas = document.createElement("canvas");
    vi.spyOn(canvas, "toBlob").mockImplementation((callback) => callback(blob));
    result.current.cropperRef.current = {
      getCanvas: () => canvas,
      getCoordinates: () => ({ left: 0, top: 200, width: 1800, height: 450 }),
    } as FixedCropperRef;
    vi.mocked(assetUpload).mockResolvedValue({ ...crop, id: "second-crop" });

    await act(async () => result.current.handleInteractionEnd());

    expect(assetUpload).toHaveBeenCalledWith(blob, {
      filename: "cropped-cover",
      parent_asset_id: "original",
    });
    expect(original.parent).toBeUndefined();
  });
});
