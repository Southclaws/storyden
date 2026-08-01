import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { NumberInput } from ".";

const meta = {
  title: "UI/Number Input",
  component: NumberInput,
  args: {
    defaultValue: "12",
    min: 0,
    max: 100,
    size: "md",
    variant: "outline",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg", "xl"],
    },
    variant: {
      control: "select",
      options: ["outline", "ghost"],
    },
  },
} satisfies Meta<typeof NumberInput>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="3" maxW="sm">
      {(["sm", "md", "lg", "xl"] as const).map((size) => (
        <NumberInput
          key={size}
          size={size}
          defaultValue="12"
          inputProps={{ "aria-label": `Number input ${size}` }}
        />
      ))}
      <NumberInput
        variant="ghost"
        defaultValue="24"
        inputProps={{ "aria-label": "Ghost number input" }}
      />
      <NumberInput
        scrubber
        defaultValue="36"
        inputProps={{ "aria-label": "Scrubbable number input" }}
      />
    </LStack>
  ),
};
