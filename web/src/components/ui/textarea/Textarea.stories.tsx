import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Textarea } from ".";

const meta = {
  title: "Components/Forms/Textarea",
  component: Textarea,
  parameters: {
    docs: {
      description: {
        component:
          "The standard multiline text control. Use it instead of `styled.textarea` for ordinary prose and instructions that need more room than an Input. Pair it with FormControl, FormLabel, helper text, and error text in forms. Use specialist composers or code editors only when the content has its own editing model.",
      },
    },
  },
  args: {
    placeholder: "Add more detail...",
    size: "md",
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
} satisfies Meta<typeof Textarea>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4" maxW="md">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Textarea
          key={size}
          size={size}
          rows={3}
          placeholder={`${size} multiline control`}
        />
      ))}
    </LStack>
  ),
};
