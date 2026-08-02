import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { styled } from "@/styled-system/jsx";

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

export const MetadataContract: Story = {
  render: () => (
    <styled.div display="flex" flexDirection="column" gap="2">
      <Text variant="metadata">Updated 12 minutes ago</Text>
      <FormFeedback>Use a clear, memorable name.</FormFeedback>
      <FormFeedback error="This slug is already in use." />
    </styled.div>
  ),
};
