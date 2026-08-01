import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormHelperText } from "../form-helper-text";
import { FormLabel } from "../form-label";
import { Input } from "../input";

import { FormControl } from ".";

const meta = {
  title: "UI/Form Control",
  component: FormControl,
} satisfies Meta<typeof FormControl>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <FormControl maxW="sm">
      <FormLabel>Community name</FormLabel>
      <Input placeholder="Storyden" />
      <FormHelperText>
        Visible in the top navigation and metadata.
      </FormHelperText>
    </FormControl>
  ),
};
