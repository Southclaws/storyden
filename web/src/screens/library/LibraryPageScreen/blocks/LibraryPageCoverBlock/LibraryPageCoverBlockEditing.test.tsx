import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { createRef } from "react";
import { beforeEach, expect, it, vi } from "vitest";

import type { Asset } from "@/api/openapi-schema";
import type { ImageCropperRef } from "@/components/ui/image-cropper";

import { LibraryPageCoverBlockEditing } from "./LibraryPageCoverBlockEditing";

let primaryImage: Asset;
let layout: "banner" | "artwork" | "image" | undefined;
beforeEach(() => {
  layout = "image";
  vi.clearAllMocks();
});

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: { getState: () => ({ setPrimaryImage: vi.fn() }) },
  }),
}));
vi.mock("../../store", () => ({
  useWatch: (selector: (state: unknown) => unknown) =>
    selector({
      draft: {
        primary_image: primaryImage,
        meta: {
          layout: { blocks: [{ type: "cover", config: { layout } }] },
        },
      },
    }),
}));

it.each([undefined, "banner", "artwork", "image"] as const)(
  "offers crop-ratio choices only for resizable Image layout (%s)",
  async (mode) => {
    layout = mode;
    primaryImage = {
      id: "photo",
      filename: "photo.png",
      path: "/api/assets/photo.png",
      mime_type: "image/png",
      width: 1600,
      height: 900,
    };
    render(
      <LibraryPageCoverBlockEditing
        cropperRef={createRef<ImageCropperRef>()}
        handleInteractionEnd={vi.fn()}
        busy={false}
      />,
    );
    const user = userEvent.setup();
    await user.click(
      screen.getByRole("button", { name: "Main image options" }),
    );
    expect(screen.getByRole("menuitem", { name: "Page layout" })).toBeVisible();
    if (mode === "image")
      expect(
        screen.getByRole("menuitem", { name: "Crop ratio" }),
      ).toBeVisible();
    else
      expect(
        screen.queryByRole("menuitem", { name: "Crop ratio" }),
      ).not.toBeInTheDocument();
  },
);
