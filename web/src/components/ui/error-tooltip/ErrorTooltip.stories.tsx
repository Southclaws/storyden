import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { ErrorTooltip } from ".";

const meta = {
  title: "UI/Error Tooltip",
  component: ErrorTooltip,
  args: {
    error: new Error("The request failed because the token expired."),
  },
} satisfies Meta<typeof ErrorTooltip>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
