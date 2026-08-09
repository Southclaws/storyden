import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormLabel } from ".";

const meta = {
  title: "Components/Forms/Form Label",
  component: FormLabel,
  args: {
    children: "Display name",
  },
} satisfies Meta<typeof FormLabel>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
