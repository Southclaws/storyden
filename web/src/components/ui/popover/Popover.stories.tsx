import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import * as Popover from ".";

const meta = {
  title: "UI/Popover",
  component: Popover.Root,
} satisfies Meta<typeof Popover.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Popover.Root>
      <Popover.Trigger asChild>
        <Button variant="outline">Open popover</Button>
      </Popover.Trigger>
      <Popover.Positioner>
        <Popover.Content>
          <Popover.Arrow>
            <Popover.ArrowTip />
          </Popover.Arrow>
          <Popover.Title>Topic settings</Popover.Title>
          <Popover.Description>
            Configure notifications, visibility, and moderation options.
          </Popover.Description>
          <Popover.CloseTrigger asChild>
            <Button size="sm" variant="ghost">
              Close
            </Button>
          </Popover.CloseTrigger>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  ),
};
