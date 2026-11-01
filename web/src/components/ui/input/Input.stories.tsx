import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack, styled } from "@/styled-system/jsx";

import { CardBox } from "../card-box";

import { Input, InputPrefix } from ".";

const meta = {
  title: "Components/Forms/Input",
  component: Input,
  parameters: {
    docs: {
      description: {
        component:
          "The standard compact text control. Use the default `sm` outline input on the control background for ordinary forms. Use `inset` only when the control is embedded in an Inset region, and `ghost` for deliberate inline-editing or composer treatments where a bordered field would add noise.",
      },
    },
  },
  args: {
    placeholder: "Search community...",
    size: "sm",
    variant: "outline",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["outline", "ghost", "inset"],
    },
  },
} satisfies Meta<typeof Input>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const InlineEditing: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Ghost fields blend into the surrounding surface while idle. Focus a field to reveal its editing outline. Use outline inputs when a persistent field boundary is needed.",
      },
    },
  },
  render: () => (
    <CardBox maxW="sm" p="4">
      <LStack gap="3">
        <Input
          aria-label="Page name"
          variant="ghost"
          defaultValue="Departure Mono"
          fontWeight="semibold"
        />
        <Input
          aria-label="Page description"
          variant="ghost"
          defaultValue="A free monospaced pixel font."
        />
        <Input
          aria-label="Page URL"
          variant="ghost"
          defaultValue="https://departuremono.com/"
        />
      </LStack>
    </CardBox>
  ),
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4" maxW="md">
      {(["outline", "ghost", "inset"] as const).map((variant) => (
        <LStack key={variant} gap="3">
          {(["sm", "md", "lg"] as const).map((size) => (
            <Input
              key={size}
              size={size}
              variant={variant}
              placeholder={`${variant} ${size}`}
            />
          ))}
        </LStack>
      ))}
    </LStack>
  ),
};

export const WithPrefix: Story = {
  render: () => (
    <styled.div display="flex" maxW="md" width="full">
      <InputPrefix>storyden.com/</InputPrefix>
      <Input borderTopLeftRadius="none" borderBottomLeftRadius="none" />
    </styled.div>
  ),
};
