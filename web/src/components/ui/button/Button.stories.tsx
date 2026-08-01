import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Button, ButtonGroup } from ".";

const meta = {
  title: "UI/Button",
  component: Button,
  args: {
    children: "Create post",
    size: "md",
    variant: "solid",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg", "xl", "2xl"],
    },
    variant: {
      control: "select",
      options: ["solid", "outline", "ghost", "link", "subtle"],
    },
  },
} satisfies Meta<typeof Button>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      {(["solid", "subtle", "outline", "ghost", "link"] as const).map(
        (variant) => (
          <HStack key={variant} gap="3" flexWrap="wrap">
            {(["xs", "sm", "md", "lg", "xl", "2xl"] as const).map((size) => (
              <Button key={size} variant={variant} size={size}>
                {variant} {size}
              </Button>
            ))}
          </HStack>
        ),
      )}
    </LStack>
  ),
};

export const States: Story = {
  render: () => (
    <HStack gap="3" flexWrap="wrap">
      <Button>Default</Button>
      <Button loading>Loading</Button>
      <Button loading loadingText="Saving...">
        Save
      </Button>
      <Button disabled>Disabled</Button>
    </HStack>
  ),
};

export const Grouped: Story = {
  render: () => (
    <ButtonGroup variant="outline" size="sm">
      <Button>Preview</Button>
      <Button>Edit</Button>
      <Button>Publish</Button>
    </ButtonGroup>
  ),
};
