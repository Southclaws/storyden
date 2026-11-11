import { act, renderHook } from "@testing-library/react";
import { expect, it, vi } from "vitest";

import { assetUpload } from "@/api/openapi-client/assets";
import type { Asset, Category } from "@/api/openapi-schema";
import type { ImageCropperSnapshot } from "@/components/ui/image-cropper";

import { useCategoryEdit } from "./useCategoryEdit";

const updateCategory = vi.fn();

vi.mock("next/navigation", () => ({
  usePathname: () => "/d/general",
  useRouter: () => ({ replace: vi.fn() }),
}));

vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<void>) => fn(),
}));

vi.mock("@/api/openapi-client/assets", () => ({ assetUpload: vi.fn() }));

vi.mock("@/lib/category/mutation", () => ({
  useCategoryMutations: () => ({ updateCategory, revalidateList: vi.fn() }),
}));

const cover: Asset = {
  id: "old-cover",
  path: "/api/assets/old-cover",
  filename: "cover.png",
  mime_type: "image/png",
  width: 1536,
  height: 384,
};

const category: Category = {
  id: "general-id",
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
  name: "General",
  slug: "general",
  description: "Community conversation",
  colour: "#8577ce",
  sort: 0,
  children: [],
  postCount: 0,
  cover_image: cover,
};

it("keeps the existing cover when saving unrelated category changes", async () => {
  const { result } = renderHook(() => useCategoryEdit({ category }));

  act(() => result.current.form.setValue("name", "Community"));
  await act(() => result.current.handlers.handleSubmit());

  expect(updateCategory).toHaveBeenCalledWith(
    "general",
    expect.objectContaining({
      name: "Community",
      cover_image_asset_id: cover.id,
    }),
  );
  expect(assetUpload).not.toHaveBeenCalled();
});

it("removes the cover only when the category is saved", async () => {
  const { result } = renderHook(() => useCategoryEdit({ category }));

  await act(() => result.current.handlers.handleImageRemove());
  expect(updateCategory).not.toHaveBeenCalled();

  await act(() => result.current.handlers.handleSubmit());
  expect(updateCategory).toHaveBeenCalledWith(
    "general",
    expect.objectContaining({ cover_image_asset_id: null }),
  );
  expect(assetUpload).not.toHaveBeenCalled();
});

it("cancels a cover removal without losing the existing cover", async () => {
  const { result } = renderHook(() => useCategoryEdit({ category }));

  await act(() => result.current.handlers.handleImageRemove());
  act(() => result.current.handlers.handleCancel());
  await act(() => result.current.handlers.handleSubmit());

  expect(updateCategory).toHaveBeenCalledWith(
    "general",
    expect.objectContaining({ cover_image_asset_id: cover.id }),
  );
});

it("saves the cropped image and attaches the uploaded asset to the category", async () => {
  const image = new Blob(["cropped image"], { type: "image/png" });
  const crop: ImageCropperSnapshot = {
    coordinates: { left: 0, top: 100, width: 800, height: 200 },
    imageSize: { width: 800, height: 600 },
    width: 1536,
    height: 384,
    getBlob: async () => image,
  };
  vi.mocked(assetUpload).mockResolvedValue({ ...cover, id: "new-cover" });
  const { result } = renderHook(() => useCategoryEdit({ category }));
  result.current.cropperRef.current = {
    getCrop: () => crop,
    getCoordinates: () => crop.coordinates,
    getCroppedImage: async () => image,
  };

  act(() => result.current.handlers.handleImageChange());
  expect(assetUpload).not.toHaveBeenCalled();

  await act(() => result.current.handlers.handleSubmit());
  expect(assetUpload).toHaveBeenCalledWith(image, expect.any(Object));
  expect(updateCategory).toHaveBeenCalledWith(
    "general",
    expect.objectContaining({ cover_image_asset_id: "new-cover" }),
  );
});

it("does not save the category if the cropped image upload fails", async () => {
  vi.mocked(assetUpload).mockRejectedValue(new Error("Upload failed"));
  const { result } = renderHook(() => useCategoryEdit({ category }));
  result.current.cropperRef.current = {
    getCrop: () => ({
      coordinates: { left: 0, top: 0, width: 800, height: 200 },
      imageSize: { width: 800, height: 600 },
      width: 1536,
      height: 384,
      getBlob: async () => new Blob(["crop"]),
    }),
    getCoordinates: () => null,
    getCroppedImage: async () => null,
  };
  act(() => result.current.handlers.handleImageChange());

  await act(async () => {
    await expect(result.current.handlers.handleSubmit()).rejects.toThrow(
      "Upload failed",
    );
  });
  expect(updateCategory).not.toHaveBeenCalled();
});

it("waits for a newly selected image before saving", async () => {
  const { result } = renderHook(() => useCategoryEdit({ category }));

  await act(() => result.current.handlers.handleImageFileChange());
  await act(() => result.current.handlers.handleSubmit());

  expect(updateCategory).not.toHaveBeenCalled();
  expect(assetUpload).not.toHaveBeenCalled();

  await act(() => result.current.handlers.handleImageRemove());
  await act(() => result.current.handlers.handleSubmit());

  expect(updateCategory).toHaveBeenCalledWith(
    "general",
    expect.objectContaining({ cover_image_asset_id: null }),
  );
});
