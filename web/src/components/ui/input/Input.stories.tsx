import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack, styled } from "@/styled-system/jsx";

import { Input, InputPrefix } from ".";

const meta = {
  title: "Components/Forms/Input",
  component: Input,
  args: {
    placeholder: "Search community...",
    size: "sm",
    variant: "outline",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["outline", "ghost", "inset"],
    },
  },
} satisfies Meta<typeof Input>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4" maxW="md">
      {(["outline", "ghost", "inset"] as const).map((variant) => (
        <LStack key={variant} gap="3">
          {(["sm", "md", "lg"] as const).map((size) => (
            <Input
              key={size}
              size={size}
              variant={variant}
              placeholder={`${variant} ${size}`}
            />
          ))}
        </LStack>
      ))}
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
