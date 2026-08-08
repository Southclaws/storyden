import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { PageHeading } from ".";

const meta = {
  title: "Components/Typography/Page Heading",
  component: PageHeading,
  args: {
    children: "System settings",
  },
} satisfies Meta<typeof PageHeading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
