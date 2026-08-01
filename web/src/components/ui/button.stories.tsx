import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Button, ButtonGroup } from "./button";

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
      <HStack gap="3" flexWrap="wrap">
        <Button variant="solid">Solid</Button>
        <Button variant="subtle">Subtle</Button>
        <Button variant="outline">Outline</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="link">Link</Button>
      </HStack>
      <HStack gap="3" flexWrap="wrap">
        <Button size="xs">Extra small</Button>
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
        <Button size="xl">Extra large</Button>
      </HStack>
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
