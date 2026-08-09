import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack } from "@/styled-system/jsx";

import { Spinner } from ".";

const meta = {
  title: "Components/Feedback/Spinner",
  component: Spinner,
  parameters: {
    docs: {
      description: {
        component:
          "Indeterminate activity for compact controls or regions whose completion cannot be measured. Preserve the control's width while loading. Use Progress for measurable work and a skeleton when loading content should retain its eventual layout.",
      },
    },
  },
  args: {
    size: "md",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Spinner>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Sizes: Story = {
  render: () => (
    <HStack gap="4">
      <Spinner size="sm" />
      <Spinner size="md" />
      <Spinner size="lg" />
    </HStack>
  ),
};
