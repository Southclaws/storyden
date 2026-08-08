import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { PinInput } from ".";

const meta = {
  title: "Components/Forms/Pin Input",
  component: PinInput,
  args: {
    children: "Verification code",
    length: 4,
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof PinInput>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4">
      {(["sm", "md", "lg"] as const).map((size) => (
        <PinInput key={size} size={size}>
          Pin {size}
        </PinInput>
      ))}
    </LStack>
  ),
};
