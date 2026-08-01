import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Checkbox } from ".";

const meta = {
  title: "UI/Checkbox",
  component: Checkbox,
  args: {
    children: "Email me replies",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Checkbox>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="3">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Checkbox key={size} size={size} defaultChecked>
          Checkbox {size}
        </Checkbox>
      ))}
    </LStack>
  ),
};

export const States: Story = {
  render: () => (
    <HStack gap="4" flexWrap="wrap">
      <Checkbox>Unchecked</Checkbox>
      <Checkbox defaultChecked>Checked</Checkbox>
      <Checkbox checked="indeterminate">Indeterminate</Checkbox>
      <Checkbox disabled>Disabled</Checkbox>
    </HStack>
  ),
};
