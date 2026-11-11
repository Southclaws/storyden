import { act, renderHook } from "@testing-library/react";
import { expect, it, vi } from "vitest";

import { assetUpload } from "@/api/openapi-client/assets";
import { categoryCreate } from "@/api/openapi-client/categories";
import type { Asset } from "@/api/openapi-schema";

import { useCategoryCreate } from "./useCategoryCreate";

vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<void>) => fn(),
}));
vi.mock("@/api/openapi-client/assets", () => ({ assetUpload: vi.fn() }));
vi.mock("@/api/openapi-client/categories", () => ({
  categoryCreate: vi.fn(),
  getCategoryListKey: () => ["/categories"],
}));

it.each([false, true])(
  "creates a category with the chosen cover, respecting removal=%s",
  async (remove) => {
    const uploaded: Asset = {
      id: "new-cover",
      path: "/api/assets/new-cover",
      filename: "cover.png",
      mime_type: "image/png",
      width: 1536,
      height: 384,
    };
    vi.mocked(assetUpload).mockResolvedValue(uploaded);
    const { result } = renderHook(() => useCategoryCreate({}));

    await act(async () => {
      for (const [name, value] of [
        ["name", "General"],
        ["slug", "general"],
        ["description", "Community conversation"],
      ] as const) {
        const field = result.current.register(name);
        await field.onChange({ target: { name, value }, type: "change" });
      }
    });

    const image = new Blob(["cropped image"], { type: "image/png" });
    result.current.cropperRef.current = {
      getCrop: () => ({
        coordinates: { left: 0, top: 0, width: 800, height: 200 },
        imageSize: { width: 800, height: 600 },
        width: 1536,
        height: 384,
        getBlob: async () => image,
      }),
      getCoordinates: () => null,
      getCroppedImage: async () => image,
    };
    act(() => result.current.handleImageChange());

    if (remove) {
      await act(() => result.current.handleImageRemove());
    }
    expect(categoryCreate).not.toHaveBeenCalled();
    expect(assetUpload).not.toHaveBeenCalled();

    await act(() => result.current.onSubmit());

    const request = vi.mocked(categoryCreate).mock.calls[0]![0]!;
    expect(request).toMatchObject({ name: "General", slug: "general" });
    if (remove) {
      expect(request.cover_image_asset_id).toBeUndefined();
      expect(assetUpload).not.toHaveBeenCalled();
    } else {
      expect(request.cover_image_asset_id).toBe("new-cover");
      expect(assetUpload).toHaveBeenCalledWith(image, expect.any(Object));
    }
  },
);
