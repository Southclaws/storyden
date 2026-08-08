import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { SectionHeading } from ".";

const meta = {
  title: "Components/Typography/Section Heading",
  component: SectionHeading,
  args: {
    children: "Client IP strategy",
  },
} satisfies Meta<typeof SectionHeading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
