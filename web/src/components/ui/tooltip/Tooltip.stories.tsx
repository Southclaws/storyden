import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box } from "@/styled-system/jsx";

import { Button } from "../button";

import * as Tooltip from ".";
import { ErrorTooltip } from "./ErrorTooltip.internal";
import { Tooltip as ClosedTooltip } from "./Tooltip";

const meta = {
  title: "Components/Overlays/Tooltip",
  component: ClosedTooltip,
  args: {
    children: <Button variant="outline">Hover for tooltip</Button>,
    content: "Create a new discussion",
  },
} satisfies Meta<typeof ClosedTooltip>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <ClosedTooltip content="Create a new discussion" openDelay={100}>
      <Button variant="outline">Hover for tooltip</Button>
    </ClosedTooltip>
  ),
};

export const Compound: Story = {
  render: () => (
    <Tooltip.Root openDelay={100}>
      <Tooltip.Trigger asChild>
        <Button variant="outline">Custom tooltip</Button>
      </Tooltip.Trigger>
      <Tooltip.Positioner>
        <Tooltip.Content>Custom composition</Tooltip.Content>
      </Tooltip.Positioner>
    </Tooltip.Root>
  ),
};

export const ViewportEdge: Story = {
  render: () => (
    <Box display="flex" justifyContent="end" minH="sm" width="full">
      <ClosedTooltip
        content="Collision-safe tooltip"
        open
        positioning={{ placement: "bottom-end" }}
      >
        <Button variant="outline">Edge trigger</Button>
      </ClosedTooltip>
    </Box>
  ),
};

export const Error: Story = {
  render: () => (
    <ErrorTooltip
      error={new globalThis.Error("The category list could not be loaded.")}
    />
  ),
};
