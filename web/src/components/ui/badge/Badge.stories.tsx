import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Badge, badgeColourPalette, badgeColours } from ".";

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
      {(["solid", "subtle", "outline"] as const).map((variant) => (
        <HStack key={variant} gap="3" flexWrap="wrap">
          {(["sm", "md", "lg"] as const).map((size) => (
            <Badge key={size} variant={variant} size={size}>
              {variant} {size}
            </Badge>
          ))}
        </HStack>
      ))}
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
