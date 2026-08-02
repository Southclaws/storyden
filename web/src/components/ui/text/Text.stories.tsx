import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { SectionHeading } from "@/components/ui/section-heading";
import { styled } from "@/styled-system/jsx";

import { Text } from ".";

const meta = {
  title: "UI/Typography/Text",
  component: Text,
  args: {
    children: "There are 4 categories available to start discussions.",
    variant: "body",
  },
  argTypes: {
    variant: {
      control: "select",
      options: ["body", "supporting", "metadata"],
    },
  },
} satisfies Meta<typeof Text>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Roles: Story = {
  render: () => (
    <styled.div
      display="flex"
      flexDirection="column"
      gap="6"
      maxW="layout.readable"
    >
      <styled.div>
        <Text variant="body">
          Robots can use connected providers and tools to work with your
          community.
        </Text>
      </styled.div>

      <styled.div>
        <SectionHeading>MCP servers</SectionHeading>
        <Text variant="supporting">
          Connect external MCP servers and add their tools to Robots.
        </Text>
      </styled.div>

      <Text as="span" variant="metadata">
        Updated 12 minutes ago
      </Text>
    </styled.div>
  ),
};

export const SemanticElements: Story = {
  render: () => (
    <styled.div
      display="flex"
      flexDirection="column"
      gap="4"
      maxW="layout.readable"
    >
      <Text as="p" variant="body">
        Discussion categories organise conversations around shared interests.
      </Text>
      <Text as="span" variant="metadata">
        4 categories
      </Text>
      <Text as="aside" variant="supporting">
        Category permissions may limit where some members can post.
      </Text>
    </styled.div>
  ),
};
