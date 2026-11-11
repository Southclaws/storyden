import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Box, LStack } from "@/styled-system/jsx";

import { PinchGestureHint } from ".";

const meta = {
  title: "Components/Feedback/Pinch Gesture Hint",
  component: PinchGestureHint,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A brief visual cue for a surface that supports two-finger zoom. Two circles spread diagonally while fading, pause briefly, then reappear apart and pinch together while fading again. Each fade starts gently and accelerates as the movement finishes. Plays once on mount and removes itself after 2 seconds; remount to replay. At every size, prefers-reduced-motion shows the expanded circles joined by a diagonal double-headed arrow, with no animation or automatic disappearance. The SVG defaults to this expanded position. Inherits currentColor and never intercepts input. This is an interaction hint, not a loading indicator or zoom control. Visually hidden text says “Use two fingers to zoom” and stays available after the animation ends. Image Cropper uses this cue on ready fixed crops and dismisses it on interaction. The Replay button and frame belong only to this preview.",
      },
    },
  },
  args: {
    size: 72,
  },
  render: function Preview(args) {
    const [playback, setPlayback] = useState(0);

    return (
      <LStack gap="4" alignItems="center">
        <Box
          display="flex"
          alignItems="center"
          justifyContent="center"
          width="64"
          maxWidth="full"
          height="32"
          background="background.inset"
          color="text.default"
          borderRadius="panel"
        >
          <PinchGestureHint key={playback} {...args} />
        </Box>
        <Button
          size="sm"
          variant="subtle"
          onClick={() => setPlayback((value) => value + 1)}
        >
          Replay animation
        </Button>
      </LStack>
    );
  },
} satisfies Meta<typeof PinchGestureHint>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Preview: Story = {};

export const Compact: Story = {
  args: { size: 48 },
};
