import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Badge, badgeColourPalette, badgeColours } from "./badge";

const meta = {
  title: "UI/Badge",
  component: Badge,
  args: {
    children: "Announcement",
    size: "md",
    variant: "subtle",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["solid", "subtle", "outline"],
    },
  },
} satisfies Meta<typeof Badge>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      <HStack gap="3" flexWrap="wrap">
        <Badge variant="solid">Solid</Badge>
        <Badge variant="subtle">Subtle</Badge>
        <Badge variant="outline">Outline</Badge>
      </HStack>
      <HStack gap="3" flexWrap="wrap">
        <Badge size="sm">Small</Badge>
        <Badge size="md">Medium</Badge>
        <Badge size="lg">Large</Badge>
      </HStack>
    </LStack>
  ),
};

export const ColourPalette: Story = {
  render: () => (
    <HStack gap="3" flexWrap="wrap">
      <Badge style={badgeColourPalette(badgeColours("#27b981"))}>Open</Badge>
      <Badge style={badgeColourPalette(badgeColours("#6d5dfc"))}>
        Featured
      </Badge>
      <Badge style={badgeColourPalette(badgeColours("#d94662"))}>Urgent</Badge>
    </HStack>
  ),
};
