import { render } from "@testing-library/react";
import { createRef } from "react";
import { FixedCropper, type FixedCropperRef } from "react-advanced-cropper";
import { expect, it, vi } from "vitest";

import type { Asset } from "@/api/openapi-schema";

import { LibraryPageCoverBlockEditing } from "./LibraryPageCoverBlockEditing";

let primaryImage: Asset;

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({
    store: { getState: () => ({ setPrimaryImage: vi.fn() }) },
  }),
}));
vi.mock("../../store", () => ({
  useWatch: (selector: (state: unknown) => unknown) =>
    selector({ draft: { primary_image: primaryImage, meta: {} } }),
}));
vi.mock("react-advanced-cropper", () => ({
  FixedCropper: vi.fn(() => null),
  ImageRestriction: { stencil: "stencil" },
}));

it("edits original pixels both before and after saving a crop", () => {
  const original: Asset = {
    id: "original",
    filename: "original.png",
    path: "/api/assets/original.png",
    mime_type: "image/png",
    width: 1800,
    height: 900,
  };
  primaryImage = original;
  const props = {
    cropperRef: createRef<FixedCropperRef>(),
    handleInteractionEnd: vi.fn(),
  };
  const { rerender } = render(<LibraryPageCoverBlockEditing {...props} />);
  expect(vi.mocked(FixedCropper).mock.lastCall?.[0].src).toBe(
    "https://api.storyden.test/api/assets/original.png",
  );

  primaryImage = {
    ...original,
    id: "crop",
    path: "/api/assets/crop.png",
    parent: original,
  };
  rerender(<LibraryPageCoverBlockEditing {...props} />);
  expect(vi.mocked(FixedCropper).mock.lastCall?.[0].src).toBe(
    "https://api.storyden.test/api/assets/original.png",
  );
});
