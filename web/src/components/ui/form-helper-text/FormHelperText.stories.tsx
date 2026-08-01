import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormHelperText } from ".";

const meta = {
  title: "UI/Form Helper Text",
  component: FormHelperText,
  args: {
    children: "Markdown links are supported.",
  },
} satisfies Meta<typeof FormHelperText>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
