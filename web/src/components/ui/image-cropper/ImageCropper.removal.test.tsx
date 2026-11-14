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
  const LoadedCropper = forwardRef(function LoadedCropper(
    { onError }: { onError?: () => void },
    ref,
  ) {
    useImperativeHandle(ref, () => ({
      getCoordinates: () => ({ left: 0, top: 0, width: 200, height: 50 }),
      getImage: () => ({ width: 200, height: 100 }),
      getCanvas: () => document.createElement("canvas"),
    }));

    return (
      <div>
        <div role="img" aria-label="Crop preview" />
        <button type="button" onClick={onError}>
          Simulate decode failure
        </button>
      </div>
    );
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

it("notifies the caller and announces when image decoding fails", async () => {
  const onLoadError = vi.fn();
  const onChange = vi.fn();
  render(
    <ImageCropper
      src="/corrupt.png"
      width={200}
      height={50}
      onLoadError={onLoadError}
      onChange={onChange}
    />,
  );

  await userEvent
    .setup()
    .click(screen.getByRole("button", { name: "Simulate decode failure" }));

  expect(onLoadError).toHaveBeenCalledOnce();
  expect(onChange).not.toHaveBeenCalled();
  expect(screen.getByRole("group", { name: "Crop area" })).toHaveAttribute(
    "aria-busy",
    "false",
  );
  expect(screen.getByRole("status")).toHaveTextContent(
    "The image could not be opened. Please choose another.",
  );
});
