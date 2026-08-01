import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import * as Tooltip from ".";

const meta = {
  title: "UI/Tooltip",
  component: Tooltip.Root,
} satisfies Meta<typeof Tooltip.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Tooltip.Root openDelay={100}>
      <Tooltip.Trigger asChild>
        <Button variant="outline">Hover for tooltip</Button>
      </Tooltip.Trigger>
      <Tooltip.Positioner>
        <Tooltip.Content>
          <Tooltip.Arrow>
            <Tooltip.ArrowTip />
          </Tooltip.Arrow>
          Create a new discussion
        </Tooltip.Content>
      </Tooltip.Positioner>
    </Tooltip.Root>
  ),
};
