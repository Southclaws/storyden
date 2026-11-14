import { act, renderHook, waitFor } from "@testing-library/react";
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

it("allows creation without a cover after an image fails to decode", async () => {
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

  await act(() => result.current.handleImageFileChange());
  expect(result.current.imageLoading).toBe(true);

  act(() => result.current.handleImageError());
  expect(result.current.imageLoading).toBe(false);
  await act(() => result.current.onSubmit());

  expect(categoryCreate).toHaveBeenCalledWith(
    expect.objectContaining({
      name: "General",
      slug: "general",
    }),
  );
  expect(
    vi.mocked(categoryCreate).mock.calls[0]?.[0]?.cover_image_asset_id,
  ).toBeUndefined();
  expect(assetUpload).not.toHaveBeenCalled();
});

it.each(["cancel", "close", "unmount"])(
  "does not create a category when the form is dismissed with %s during upload",
  async (dismiss) => {
    const upload = Promise.withResolvers<Asset>();
    vi.mocked(assetUpload).mockReturnValue(upload.promise);
    const onClose = vi.fn();
    const { result, rerender, unmount } = renderHook(
      ({ isOpen }) => useCategoryCreate({ onClose, isOpen }),
      { initialProps: { isOpen: true } },
    );

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

    let submission: Promise<void>;
    await act(async () => {
      submission = result.current.onSubmit();
    });
    await waitFor(() => expect(assetUpload).toHaveBeenCalledOnce());

    if (dismiss === "cancel") {
      act(() => result.current.handleCancel());
      expect(onClose).toHaveBeenCalledOnce();
    } else if (dismiss === "close") {
      rerender({ isOpen: false });
    } else {
      unmount();
    }
    onClose.mockClear();

    await act(async () => {
      upload.resolve({
        id: "new-cover",
        path: "/api/assets/new-cover",
        filename: "cover.png",
        mime_type: "image/png",
        width: 1536,
        height: 384,
      });
      await submission;
    });

    expect(categoryCreate).not.toHaveBeenCalled();
    expect(onClose).not.toHaveBeenCalled();
  },
);
