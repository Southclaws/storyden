import { fireEvent, render } from "@testing-library/react";
import { withNuqsTestingAdapter } from "nuqs/adapters/testing";
import { describe, expect, test, vi } from "vitest";

import type { Asset } from "@/api/openapi-schema";

import { AssetLightbox } from "./AssetLightbox";

const asset: Asset = {
  id: "mountains",
  filename: "mountains.png",
  path: "/assets/mountains.png",
  mime_type: "image/png",
  width: 1200,
  height: 600,
};

describe("AssetLightbox", () => {
  test("Escape closes only the visible lightbox", () => {
    const onClose = vi.fn();
    const onHiddenClose = vi.fn();
    render(
      <>
        <AssetLightbox asset={asset} present onClose={onClose} />
        <AssetLightbox asset={asset} present={false} onClose={onHiddenClose} />
      </>,
      { wrapper: withNuqsTestingAdapter() },
    );

    fireEvent.keyDown(document, { key: "Enter" });
    expect(onClose).not.toHaveBeenCalled();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalledOnce();
    expect(onHiddenClose).not.toHaveBeenCalled();
  });

  test("stops handling Escape when hidden or unmounted", () => {
    const onClose = vi.fn();
    const { rerender, unmount } = render(
      <AssetLightbox asset={asset} present onClose={onClose} />,
      { wrapper: withNuqsTestingAdapter() },
    );

    rerender(<AssetLightbox asset={asset} present={false} onClose={onClose} />);
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).not.toHaveBeenCalled();

    rerender(<AssetLightbox asset={asset} present onClose={onClose} />);
    unmount();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).not.toHaveBeenCalled();
  });

  test("uses the current close callback after rerendering", () => {
    const originalClose = vi.fn();
    const currentClose = vi.fn();
    const { rerender } = render(
      <AssetLightbox asset={asset} present onClose={originalClose} />,
      { wrapper: withNuqsTestingAdapter() },
    );
    rerender(<AssetLightbox asset={asset} present onClose={currentClose} />);

    fireEvent.keyDown(document, { key: "Escape" });
    expect(originalClose).not.toHaveBeenCalled();
    expect(currentClose).toHaveBeenCalledOnce();
  });
});
