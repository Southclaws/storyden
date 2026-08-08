import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { Button, ButtonGroup } from ".";

const meta = {
  title: "Components/Actions/Button",
  component: Button,
  args: {
    children: "Create post",
    size: "sm",
    variant: "subtle",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["solid", "outline", "ghost", "subtle", "plain"],
    },
  },
} satisfies Meta<typeof Button>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      {(["solid", "subtle", "outline", "ghost", "plain"] as const).map(
        (variant) => (
          <HStack key={variant} gap="3" flexWrap="wrap">
            {(["sm", "md", "lg"] as const).map((size) => (
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
      <Button loading>Default</Button>
      <Button loading loadingText="Saving...">
        Save
      </Button>
      <Button disabled>Disabled</Button>
    </HStack>
  ),
};

export const DisabledVariants: Story = {
  render: () => (
    <HStack gap="3" flexWrap="wrap">
      {(["solid", "subtle", "outline", "ghost", "plain"] as const).map(
        (variant) => (
          <Button key={variant} variant={variant} disabled>
            {variant}
          </Button>
        ),
      )}
    </HStack>
  ),
};

export const Grouped: Story = {
  render: () => (
    <ButtonGroup variant="outline" size="md">
      <Button>Preview</Button>
      <Button>Edit</Button>
      <Button>Publish</Button>
    </ButtonGroup>
  ),
};
