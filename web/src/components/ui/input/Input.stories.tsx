import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack, styled } from "@/styled-system/jsx";

import { Input, InputPrefix } from ".";

const meta = {
  title: "UI/Input",
  component: Input,
  args: {
    placeholder: "Search community...",
    size: "md",
    variant: "outline",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["2xs", "xs", "sm", "md", "lg", "xl", "2xl"],
    },
    variant: {
      control: "select",
      options: ["outline", "ghost"],
    },
  },
} satisfies Meta<typeof Input>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3" maxW="md">
      <Input size="xs" placeholder="Extra small" />
      <Input size="sm" placeholder="Small" />
      <Input size="md" placeholder="Medium" />
      <Input size="lg" placeholder="Large" />
      <Input size="xl" placeholder="Extra large" />
    </LStack>
  ),
};

export const WithPrefix: Story = {
  render: () => (
    <styled.div display="flex" maxW="md" width="full">
      <InputPrefix>storyden.com/</InputPrefix>
      <Input borderTopLeftRadius="none" borderBottomLeftRadius="none" />
    </styled.div>
  ),
};
