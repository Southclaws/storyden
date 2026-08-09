import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box } from "@/styled-system/jsx";

import { Button } from "../button";

import * as Popover from ".";
import { Popover as ClosedPopover } from "./Popover";

const meta = {
  title: "Components/Overlays/Popover",
  component: ClosedPopover,
  parameters: {
    docs: {
      description: {
        component:
          "Temporary supporting content anchored to a trigger. Prefer the closed Popover so portalling, placement, collision handling, and viewport fitting use Storyden defaults. Use Menu for a command list, Tooltip for a short label, and Dialog for work that requires focused interruption.",
      },
    },
  },
  args: {
    trigger: <Button variant="outline">Open popover</Button>,
    children: <Popover.Title>Topic settings</Popover.Title>,
  },
} satisfies Meta<typeof ClosedPopover>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <ClosedPopover trigger={<Button variant="outline">Open popover</Button>}>
      <Popover.Title>Topic settings</Popover.Title>
      <Popover.Description>
        Configure notifications, visibility, and moderation options.
      </Popover.Description>
      <Popover.CloseTrigger asChild>
        <Button size="sm" variant="ghost">
          Close
        </Button>
      </Popover.CloseTrigger>
    </ClosedPopover>
  ),
};

export const Compound: Story = {
  render: () => (
    <Popover.Root>
      <Popover.Trigger asChild>
        <Button variant="outline">Custom popover</Button>
      </Popover.Trigger>
      <Popover.Positioner>
        <Popover.Content>Custom composition</Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  ),
};

export const ViewportEdge: Story = {
  render: () => (
    <Box display="flex" justifyContent="end" minH="sm" width="full">
      <ClosedPopover
        defaultOpen
        positioning={{ placement: "bottom-end" }}
        trigger={<Button variant="outline">Edge trigger</Button>}
      >
        <Popover.Title>Collision-safe popover</Popover.Title>
        <Popover.Description>
          This content remains inside the viewport.
        </Popover.Description>
      </ClosedPopover>
    </Box>
  ),
};
