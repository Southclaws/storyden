import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { assetUpload } from "@/api/openapi-client/assets";
import type { Asset } from "@/api/openapi-schema";
import type { ImageCropperSnapshot } from "@/components/ui/image-cropper";

import { useLibraryPageCoverBlock } from "./useLibraryPageCoverBlock";

let primaryImage: Asset | undefined;
const setPrimaryImage = vi.fn();
const overwriteBlock = vi.fn();

vi.mock("@/api/openapi-client/assets", () => ({ assetUpload: vi.fn() }));
vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<unknown>) => fn(),
}));
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    setImageSaving: vi.fn(),
    store: {
      getState: () => ({
        draft: {
          primary_image: primaryImage,
          meta: {
            layout: {
              blocks: [
                {
                  type: "cover",
                  config: { layout: "image", alt: "An illustration" },
                },
              ],
            },
          },
        },
        setPrimaryImage,
        overwriteBlock,
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

afterEach(() => vi.useRealTimers());

beforeEach(() => {
  vi.useFakeTimers();
  primaryImage = undefined;
  setPrimaryImage.mockReset();
  overwriteBlock.mockReset();
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
      const coordinates = { left: 0, top: 100, width: 1800, height: 450 };
      result.current.cropperRef.current = cropper(() => ({
        coordinates,
        imageSize: { width: 1800, height: 900 },
        width: 1800,
        height: 450,
        getBlob: async () => blob,
      }));
      vi.mocked(assetUpload).mockResolvedValue({
        ...crop,
        parent: replacement,
      });

      await act(async () => {
        result.current.handleInteractionEnd();
        await vi.advanceTimersByTimeAsync(500);
      });

      expect(assetUpload).toHaveBeenCalledWith(blob, {
        filename: "cropped-cover",
        parent_asset_id: "replacement",
      });
      expect(overwriteBlock).toHaveBeenCalledWith({
        type: "cover",
        config: {
          layout: "image",
          alt: "An illustration",
          crop: { left: 0, top: 1 / 9, width: 1, height: 0.5 },
        },
      });
      expect(setPrimaryImage).toHaveBeenCalledWith({
        isReplacement: false,
        asset: { ...crop, parent: replacement, width: 1800, height: 450 },
        config: { left: 0, top: 100 },
      });
    },
  );

  it("keeps repeated crops as siblings of the original", async () => {
    primaryImage = crop;
    const { result } = renderHook(() => useLibraryPageCoverBlock());
    const blob = new Blob(["second crop"], { type: "image/png" });
    result.current.cropperRef.current = cropper(() => ({
      coordinates: { left: 0, top: 200, width: 1800, height: 450 },
      imageSize: { width: 1800, height: 900 },
      width: 1800,
      height: 450,
      getBlob: async () => blob,
    }));
    vi.mocked(assetUpload).mockResolvedValue({ ...crop, id: "second-crop" });

    await act(async () => {
      result.current.handleInteractionEnd();
      await vi.advanceTimersByTimeAsync(500);
    });

    expect(assetUpload).toHaveBeenCalledWith(blob, {
      filename: "cropped-cover",
      parent_asset_id: "original",
    });
    expect(original.parent).toBeUndefined();
  });

  it("coalesces rapid adjustments and saves the latest crop even after the ref unmounts", async () => {
    primaryImage = original;
    const { result, unmount } = renderHook(() => useLibraryPageCoverBlock());
    const blob = new Blob(["latest pixels"], { type: "image/png" });
    let left = 0;
    result.current.cropperRef.current = cropper(() => ({
      coordinates: { left, top: 100, width: 900, height: 450 },
      imageSize: { width: 1800, height: 900 },
      width: 900,
      height: 450,
      getBlob: async () => blob,
    }));
    vi.mocked(assetUpload).mockClear();
    vi.mocked(assetUpload).mockResolvedValue({ ...crop, id: "latest-crop" });
    act(() => result.current.handleInteractionEnd());
    left = 200;
    act(() => result.current.handleInteractionEnd());
    await act(async () => {
      await vi.advanceTimersByTimeAsync(499);
    });
    expect(assetUpload).not.toHaveBeenCalled();
    result.current.cropperRef.current = null;
    unmount();
    await act(async () => {
      await vi.advanceTimersByTimeAsync(1);
    });
    expect(assetUpload).toHaveBeenCalledTimes(1);
    expect(setPrimaryImage).toHaveBeenCalledWith(
      expect.objectContaining({ config: { left: 200, top: 100 } }),
    );
  });
});

function cropper(getCrop: () => ImageCropperSnapshot) {
  return {
    getCrop,
    getCoordinates: () => getCrop().coordinates,
    getCroppedImage: () => getCrop().getBlob(),
  };
}
