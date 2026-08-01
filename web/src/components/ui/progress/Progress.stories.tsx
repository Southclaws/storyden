import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { ProgressCircle, ProgressHorizontal } from ".";

const meta = {
  title: "UI/Progress",
  component: ProgressHorizontal,
  args: {
    value: 64,
    size: "md",
    children: "Upload progress",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof ProgressHorizontal>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Horizontal: Story = {
  render: (args) => (
    <ProgressHorizontal {...args} maxW="md">
      {args.children}
    </ProgressHorizontal>
  ),
};

export const Circle: Story = {
  render: (args) => <ProgressCircle {...args}>{args.children}</ProgressCircle>,
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="5" maxW="md">
      {(["sm", "md", "lg"] as const).map((size) => (
        <ProgressHorizontal key={size} size={size} value={64}>
          Progress {size}
        </ProgressHorizontal>
      ))}
      <HStack gap="5">
        {(["sm", "md", "lg"] as const).map((size) => (
          <ProgressCircle key={size} size={size} value={64}>
            {size}
          </ProgressCircle>
        ))}
      </HStack>
    </LStack>
  ),
};
