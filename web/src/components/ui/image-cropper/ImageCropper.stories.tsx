import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import * as Menu from "../menu";

import { ImageCropper } from ".";

const source = `data:image/svg+xml,${encodeURIComponent(
  '<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900"><rect width="1600" height="900" fill="#155e75"/><circle cx="800" cy="450" r="280" fill="#fbbf24"/><rect x="700" y="350" width="200" height="200" fill="#164e63"/></svg>',
)}`;

const meta = {
  title: "Components/Forms/Image Cropper",
  component: ImageCropper,
  parameters: {
    docs: {
      description: {
        component:
          "Inline image cropping with a floating toolbar. Use it to edit pixels, not to define page layouts: the caller owns the surrounding frame, original asset, upload, debounced saving and persistence. Defaults retain the index cover interaction: a fixed viewport, exact width × height PNG export, collapsed tools and automatic cropping after upload. Opt into cropMode=resizable for drag handles and a Crop ratio submenu, toolbar=expanded for persistent controls, and maxOutputSize for proportional exports bounded in both dimensions. cropOnUpload=false leaves the uploaded original intact and requires the caller to update src. Keep frame dimensions tied to the original while editing, not the newly cropped output. getCrop() captures coordinates, original dimensions and an export that remains usable after unmount; getCoordinates()/getCroppedImage() remain available for existing callers. defaultCoordinates accepts pixel coordinates or a function of original dimensions. Children add caller-specific menu items. Fixed croppers support pinch and wheel zoom, with fine 2.5% wheel steps for trackpad precision; plus and minus keys provide keyboard zoom. Arrow keys move the image; Shift makes larger adjustments. Loading, saving, success and errors are announced. File Upload selects without cropping; AssetUploadEditor adds asset-specific behavior. Do not use this for ordinary file uploads or image viewing.",
      },
    },
  },
  args: {
    width: 1920,
    height: 480,
  },
} satisfies Meta<typeof ImageCropper>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Empty: Story = {};

export const Cover: Story = {
  args: {
    src: source,
  },
};

export const Removable: Story = {
  render: function RemovableCropper(args) {
    const [src, setSrc] = useState<string | undefined>(source);

    return (
      <ImageCropper
        {...args}
        src={src}
        onRemove={async () => setSrc(undefined)}
      />
    );
  },
  parameters: {
    docs: {
      description: {
        story:
          "Hover over the edit icon or activate it with the keyboard to reveal replacement and removal actions. Fixed crops zoom with pinch, wheel, or plus/minus keys. The same toolbar shows a spinner until the source image is ready in the cropper and while saving changes.",
      },
    },
  },
};

export const RestoredCrop: Story = {
  args: {
    src: source,
    defaultCoordinates: { left: 0, top: 400, width: 1600, height: 400 },
  },
};

export const Square: Story = {
  args: { src: source, width: 256, height: 256 },
};

export const Portrait: Story = {
  args: { src: source, width: 320, height: 480 },
};

export const Disabled: Story = {
  args: { src: source, disabled: true },
};

export const Saving: Story = {
  args: { src: source, saving: true },
};

export const Saved: Story = {
  args: { src: source, saved: true },
  parameters: {
    docs: {
      description: {
        story:
          "A successful save briefly shows a checkmark, then fades away without moving the image or surrounding content.",
      },
    },
  },
};

export const Resizable: Story = {
  args: {
    src: source,
    width: 1600,
    height: 900,
    cropMode: "resizable",
    toolbar: "expanded",
    maxOutputSize: { width: 1920, height: 1920 },
  },
  parameters: {
    docs: {
      description: {
        story:
          "Free selection with ratio choices in the floating menu. Changing the crop leaves the editing workspace stable; the caller persists onChange using getCrop().",
      },
    },
  },
};

export const ExpandedArtwork: Story = {
  args: {
    src: source,
    width: 1024,
    height: 1024,
    toolbar: "expanded",
    style: { maxWidth: 320, marginInline: "auto" },
  },
  parameters: {
    docs: {
      description: {
        story:
          "A square fixed viewport in a compact frame. It fills the final layout and offers no resizable selection or ratio choices.",
      },
    },
  },
};

export const RestoredProportionalCrop: Story = {
  args: {
    src: source,
    width: 1600,
    height: 900,
    cropMode: "resizable",
    toolbar: "expanded",
    maxOutputSize: { width: 1920, height: 1920 },
    defaultCoordinates: (size) => ({
      left: size.width * 0.25,
      top: size.height * 0.25,
      width: size.width * 0.5,
      height: size.height * 0.5,
    }),
  },
};

export const AdditionalMenuItems: Story = {
  args: {
    src: source,
    toolbar: "expanded",
    children: <Menu.Item value="caption">Edit image caption</Menu.Item>,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Extend the floating menu with domain-specific controls. The shared cropper continues to own replacement, removal, ratio choices and export.",
      },
    },
  },
};

export const ExtremePortrait: Story = {
  args: {
    src: `data:image/svg+xml,${encodeURIComponent('<svg xmlns="http://www.w3.org/2000/svg" width="1" height="1000"><rect width="1" height="1000" fill="#155e75"/></svg>')}`,
    width: 1,
    height: 1000,
    cropMode: "resizable",
    toolbar: "expanded",
    style: { maxHeight: "min(70svh,40rem)", minHeight: 36 },
    maxOutputSize: { width: 1920, height: 1920 },
  },
};
