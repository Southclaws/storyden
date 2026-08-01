import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Text } from ".";

const meta = {
  title: "UI/Text",
  component: Text,
  args: {
    children: "Community interfaces should be calm, legible, and direct.",
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg", "xl", "2xl"],
    },
    variant: {
      control: "select",
      options: [undefined, "heading"],
    },
  },
} satisfies Meta<typeof Text>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="2">
      {(["xs", "sm", "md", "lg", "xl", "2xl"] as const).map((size) => (
        <Text key={size} size={size}>
          Text {size}
        </Text>
      ))}
    </LStack>
  ),
};
