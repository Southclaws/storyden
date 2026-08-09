import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import { showUndoToast } from ".";

const meta = {
  title: "Components/Feedback/Undo Toast",
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Button
      onClick={() =>
        showUndoToast({
          message: "Discussion deleted",
          duration: 5000,
          onUndo: () => undefined,
          onComplete: () => undefined,
        })
      }
    >
      Show undo toast
    </Button>
  ),
};
