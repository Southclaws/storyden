import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import * as RadioGroup from ".";

const options = [
  { label: "Public", value: "public" },
  { label: "Members", value: "members" },
  { label: "Private", value: "private" },
];

const meta = {
  title: "UI/Radio Group",
  component: RadioGroup.Root,
  args: {
    size: "md",
    defaultValue: "public",
    orientation: "vertical",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    orientation: {
      control: "select",
      options: ["vertical", "horizontal"],
    },
  },
} satisfies Meta<typeof RadioGroup.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: RadioGroup.RootProps) {
  return (
    <RadioGroup.Root {...props}>
      {options.map((item) => (
        <RadioGroup.Item key={item.value} value={item.value}>
          <RadioGroup.ItemControl />
          <RadioGroup.ItemText>{item.label}</RadioGroup.ItemText>
          <RadioGroup.ItemHiddenInput />
        </RadioGroup.Item>
      ))}
    </RadioGroup.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Example key={size} size={size} defaultValue="public" />
      ))}
      <HStack gap="6">
        <Example size="md" orientation="horizontal" defaultValue="members" />
      </HStack>
    </LStack>
  ),
};
