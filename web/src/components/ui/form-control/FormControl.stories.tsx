import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { FormHelperText } from "../form-helper-text";
import { FormLabel } from "../form-label";
import { Input } from "../input";

import { FormControl } from ".";

const meta = {
  title: "Components/Forms/Form Control",
  component: FormControl,
  parameters: {
    docs: {
      description: {
        component:
          "The canonical vertical field composition: label, control, then supporting or error feedback. Keep related fields in a compact stack and use the component's built-in spacing rather than adding local margins between its parts.",
      },
    },
  },
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
