import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormErrorText } from ".";

const meta = {
  title: "Components/Forms/Form Error Text",
  component: FormErrorText,
  args: {
    children: "A community name is required.",
  },
} satisfies Meta<typeof FormErrorText>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
