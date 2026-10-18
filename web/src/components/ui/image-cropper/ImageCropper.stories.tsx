import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

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
          "Choose and position an image inside a fixed crop. Focus the crop area and use arrow keys to reposition; hold Shift for larger adjustments. Loading, saving, success and image errors are announced. Saving fades the image and disables crop, replacement and removal actions. The caller can retain the selected file with onFileChange before the local preview opens. getCroppedImage() exports a PNG; getCoordinates() and defaultCoordinates preserve the crop between editing sessions. onChange fires after a local image loads and when a crop interaction ends. The caller owns saving, upload and persistence. File Upload selects files without cropping; AssetUploadEditor includes asset-specific upload behavior. Do not use this component for ordinary file uploads or image viewing.",
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
          "Hover over the edit icon or activate it with the keyboard to reveal replacement and removal actions. The same toolbar shows a spinner until the source image is ready in the cropper and while saving changes.",
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
