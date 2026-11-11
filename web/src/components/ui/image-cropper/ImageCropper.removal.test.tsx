import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { createRef, useState } from "react";
import { expect, it, vi } from "vitest";

import { ImageCropper, type ImageCropperRef } from ".";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

vi.mock("react-advanced-cropper", async (importOriginal) => {
  const original =
    await importOriginal<typeof import("react-advanced-cropper")>();
  const { forwardRef, useImperativeHandle } = await import("react");

  // The engine can retain a loaded image after its source is cleared.
  const LoadedCropper = forwardRef(function LoadedCropper(_props, ref) {
    useImperativeHandle(ref, () => ({
      getCoordinates: () => ({ left: 0, top: 0, width: 200, height: 50 }),
      getImage: () => ({ width: 200, height: 100 }),
      getCanvas: () => document.createElement("canvas"),
    }));

    return <div role="img" aria-label="Crop preview" />;
  });

  return { ...original, Cropper: LoadedCropper, FixedCropper: LoadedCropper };
});

it.each(["fixed", "resizable"] as const)(
  "clears the %s preview and export buffer when its image is removed",
  async (cropMode) => {
    const ref = createRef<ImageCropperRef>();

    function RemovableCover() {
      const [src, setSrc] = useState<string | undefined>("/photo.png");

      return (
        <ImageCropper
          ref={ref}
          src={src}
          width={200}
          height={50}
          cropMode={cropMode}
          onRemove={async () => setSrc(undefined)}
        />
      );
    }

    render(<RemovableCover />);
    expect(ref.current?.getCrop()).not.toBeNull();
    expect(
      screen.getByRole("img", { name: "Crop preview" }),
    ).toBeInTheDocument();

    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: "Image tools" }));
    await user.click(screen.getByRole("button", { name: "Remove image" }));

    expect(
      screen.queryByRole("img", { name: "Crop preview" }),
    ).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Upload image" })).toBeVisible();
    expect(ref.current?.getCrop()).toBeNull();
    expect(ref.current?.getCoordinates()).toBeNull();
    await expect(ref.current?.getCroppedImage()).resolves.toBeNull();
  },
);
