import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "@/components/ui/button";
import { CloseIcon } from "@/components/ui/icons/Close";
import { Text } from "@/components/ui/text";

import * as FloatingPanel from ".";

const meta = {
  title: "Components/Overlays/Floating Panel",
  component: FloatingPanel.Root,
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "A persistent, non-modal work surface that can be moved and resized while the underlying page remains usable. Use it for tools that must follow someone across the product. Use Popover for short anchored content and Dialog for focused interruption.",
      },
    },
  },
} satisfies Meta<typeof FloatingPanel.Root>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <FloatingPanel.Root
      defaultOpen
      defaultPosition={{ x: 24, y: 24 }}
      defaultSize={{ width: 440, height: 320 }}
      minSize={{ width: 320, height: 240 }}
    >
      <FloatingPanel.Positioner>
        <FloatingPanel.Content>
          <FloatingPanel.DragTrigger>
            <FloatingPanel.Header>
              <FloatingPanel.Title>Theme editor</FloatingPanel.Title>
              <FloatingPanel.Control>
                <FloatingPanel.CloseTrigger asChild>
                  <Button aria-label="Close panel" size="sm" variant="ghost">
                    <CloseIcon />
                  </Button>
                </FloatingPanel.CloseTrigger>
              </FloatingPanel.Control>
            </FloatingPanel.Header>
          </FloatingPanel.DragTrigger>
          <FloatingPanel.Body padding="4">
            <Text>
              Floating panels keep the current page interactive while a
              persistent tool remains available.
            </Text>
          </FloatingPanel.Body>
          <FloatingPanel.ResizeTrigger axis="e" />
          <FloatingPanel.ResizeTrigger axis="s" />
          <FloatingPanel.ResizeTrigger axis="se" />
        </FloatingPanel.Content>
      </FloatingPanel.Positioner>
    </FloatingPanel.Root>
  ),
};
