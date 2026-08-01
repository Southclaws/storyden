import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack } from "@/styled-system/jsx";

import { Spinner } from ".";

const meta = {
  title: "UI/Spinner",
  component: Spinner,
  args: {
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Spinner>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <HStack gap="4">
      <Spinner size="sm" />
      <Spinner size="md" />
      <Spinner size="lg" />
    </HStack>
  ),
};
