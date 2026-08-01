import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormFeedback } from ".";

const meta = {
  title: "UI/Form Feedback",
  component: FormFeedback,
  args: {
    children: "Use a clear, memorable name.",
  },
} satisfies Meta<typeof FormFeedback>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Helper: Story = {};

export const Error: Story = {
  args: {
    error: "This slug is already in use.",
  },
};
