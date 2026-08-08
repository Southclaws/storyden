import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { LinkButton } from ".";

const meta = {
  title: "Components/Actions/Link Button",
  component: LinkButton,
  args: {
    href: "/",
    children: "Open discussion",
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
      options: ["solid", "outline", "ghost", "subtle"],
    },
  },
} satisfies Meta<typeof LinkButton>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      <HStack gap="3" flexWrap="wrap">
        <LinkButton href="/" variant="solid">
          Solid
        </LinkButton>
        <LinkButton href="/" variant="subtle">
          Subtle
        </LinkButton>
        <LinkButton href="/" variant="outline">
          Outline
        </LinkButton>
        <LinkButton href="/" variant="ghost">
          Ghost
        </LinkButton>
      </HStack>
      <LinkButton href="https://www.storyden.org" variant="outline">
        External link
      </LinkButton>
    </LStack>
  ),
};
