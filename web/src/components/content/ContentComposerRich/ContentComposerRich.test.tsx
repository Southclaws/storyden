import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import type { Asset } from "@/api/openapi-schema";

import { ContentComposerRich } from "./ContentComposerRich";

const upload = vi.hoisted(() => vi.fn());

vi.mock("../useImageUpload", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../useImageUpload")>();

  return {
    ...actual,
    useImageUpload: () => ({
      upload: vi.fn(),
      uploadWithProgress: upload,
    }),
  };
});

describe("ContentComposerRich", () => {
  it("renders the image upload node view while pending and replaces it on success", async () => {
    const user = userEvent.setup();
    const onAssetUpload = vi.fn();
    const uploadAsset = { path: "uploads/sunset.png" } as Asset;
    let completeUpload: (asset: Asset) => void;

    upload.mockImplementation(
      (_file: File, onProgress: (progress: number) => void) =>
        new Promise<Asset>((resolve) => {
          completeUpload = (asset) => {
            onProgress(100);
            resolve(asset);
          };
        }),
    );
    class TestURL extends URL {
      static override createObjectURL = vi.fn(() => "blob:upload-preview");
      static override revokeObjectURL = vi.fn();
    }
    vi.stubGlobal("URL", TestURL);

    const { container } = render(
      <ContentComposerRich onAssetUpload={onAssetUpload} />,
    );

    const fileInput = await waitFor(() => {
      const input =
        container.querySelector<HTMLInputElement>('input[type="file"]');
      if (!input) {
        throw new Error("Image upload input is unavailable");
      }
      return input;
    });
    await user.upload(
      fileInput,
      new File(["image"], "sunset.png", { type: "image/png" }),
    );

    const preview = await screen.findByAltText("sunset.png");
    expect(preview).toHaveAttribute("data-uploading", "true");
    expect(preview.closest("[data-node-view-wrapper]")).not.toBeNull();

    completeUpload!(uploadAsset);

    await waitFor(() => {
      expect(preview).not.toHaveAttribute("data-uploading");
      expect(preview).toHaveAttribute(
        "src",
        expect.stringContaining(uploadAsset.path),
      );
    });
    expect(onAssetUpload).toHaveBeenCalledWith(uploadAsset);
  });
});
