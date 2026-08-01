import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import { Slot } from ".";

const meta = {
  title: "UI/Slot",
  component: Slot,
} satisfies Meta<typeof Slot>;

export default meta;

type Story = StoryObj<typeof meta>;

export const ClonesSingleChild: Story = {
  render: () => (
    <Slot aria-label="Slot-provided label">
      <Button variant="outline">Receives props from Slot</Button>
    </Slot>
  ),
};
