import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { LStack } from "@/styled-system/jsx";

import { Button } from "../button";
import * as Menu from "../menu";
import { Text } from "../text";

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
          "Inline image cropping with a floating toolbar. Use it to edit pixels, not to define page layouts: the caller owns the surrounding frame, original asset, upload, debounced saving and persistence. Defaults retain the index cover interaction: a fixed viewport, exact width × height PNG export, collapsed tools and automatic cropping after upload. Image tools is a click toggle; hovering never expands the strip. Opt into cropMode=resizable for drag handles and a Crop ratio submenu, and maxOutputSize for proportional exports bounded in both dimensions. cropOnUpload=false leaves the uploaded original intact and requires the caller to update src. Keep frame dimensions tied to the original while editing, not the newly cropped output. getCrop() captures coordinates, original dimensions and an export that remains usable after unmount; getCoordinates()/getCroppedImage() remain available for existing callers. defaultCoordinates accepts pixel coordinates or a function of original dimensions. safeArea adds a centred composition guide with an accessible description and no text overlay, without changing the crop or exported pixels; use it for a narrower social-card preview inside a fixed banner. The caller supplies visible guidance outside the image; the cropper keeps keyboard instructions available to assistive technology. Editable images use a move cursor. Below the md viewport breakpoint (768px), ready fixed crops show a one-shot Pinch Gesture Hint; pointer, wheel, or keyboard input dismisses it. Reduced motion shows a static cue until interaction. Set showGestureHint=false to suppress it. Children add caller-specific menu items. Fixed croppers support pinch and wheel zoom, with fine 2.5% wheel steps for trackpad precision; plus and minus keys provide keyboard zoom. Arrow keys move the image; Shift makes larger adjustments. Loading, saving, success and errors are announced. onLoadError notifies the caller when decoding or loading fails, so enclosing forms can clear pending image state without changing the saved asset. File Upload selects files without cropping. Do not use this for ordinary file uploads or image viewing.",
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
          "Activate Image tools to reveal replacement and removal actions; activate it again to close the strip. Hovering does not open it. Clearing the image removes the preview and export buffer; getCrop() and getCroppedImage() return null until another image is ready. Fixed crops zoom with pinch, wheel, or plus/minus keys. The same toolbar shows a spinner until the source image is ready in the cropper and while saving changes.",
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

export const Artwork: Story = {
  args: {
    src: source,
    width: 1024,
    height: 1024,
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
    children: <Menu.Item value="caption">Edit image caption</Menu.Item>,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Open Image tools, then Image options to find domain-specific controls. The shared cropper continues to own replacement, removal, ratio choices and export.",
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
    style: { maxHeight: "min(70svh,40rem)", minHeight: 36 },
    maxOutputSize: { width: 1920, height: 1920 },
  },
};

export const SocialPreviewSafeArea: Story = {
  args: {
    src: source,
    safeArea: { width: 1200, height: 630, label: "Social preview" },
  },
  parameters: {
    docs: {
      description: {
        story:
          "The full 4:1 banner remains the crop. A subtle centre guide marks a common 1200:630 social-card shape; keep important details here. The guide has no text overlay and never appears in the exported PNG. Explain its purpose at the call site, outside the image. Other crops on sharing platforms may vary.",
      },
    },
  },
};

export const MobileCover: Story = {
  parameters: {
    layout: "centered",
    docs: {
      description: {
        story:
          "A mobile-width category cover with a fixed 4:1 crop and social-preview guide. Below 768px viewport width, the zoom cue starts after the image is ready, plays once, and dismisses immediately on pointer, wheel, or keyboard input. Reduced motion keeps a static expanded gesture with an arrow until interaction. A translucent black scrim keeps the gesture visible over bright images, fading away with the animation. Reduced motion keeps the scrim with the static cue until interaction. The cue and scrim are hidden at 768px and above, regardless of the cropper’s width. They never intercept input or appear in the exported image. Removal and replacement use the regular toolbar. Replay preview restores the sample image and remounts the editor for review.",
      },
    },
  },
  render: function MobileCoverPreview(args) {
    const [playback, setPlayback] = useState(0);
    const [src, setSrc] = useState<string | undefined>(source);

    return (
      <LStack width="full" maxWidth="96" minWidth="0" gap="3">
        <Text fontWeight="semibold">Cover image</Text>
        <ImageCropper
          {...args}
          key={playback}
          aria-label="Category cover"
          width={1536}
          height={384}
          src={src}
          safeArea={{ width: 1200, height: 630, label: "Social preview" }}
          onRemove={async () => setSrc(undefined)}
        />
        <Text variant="supporting">
          Keep important details inside the dashed area for social previews.
        </Text>
        <Button
          size="sm"
          variant="subtle"
          onClick={() => {
            setSrc(source);
            setPlayback((value) => value + 1);
          }}
        >
          Replay preview
        </Button>
      </LStack>
    );
  },
};
