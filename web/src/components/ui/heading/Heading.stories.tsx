import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Heading } from ".";

const meta = {
  title: "UI/Heading",
  component: Heading,
  args: {
    children: "Discussion categories",
    size: "lg",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg", "xl", "2xl"],
    },
  },
} satisfies Meta<typeof Heading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="2">
      {(["xs", "sm", "md", "lg", "xl", "2xl"] as const).map((size) => (
        <Heading key={size} size={size}>
          Heading {size}
        </Heading>
      ))}
    </LStack>
  ),
};
