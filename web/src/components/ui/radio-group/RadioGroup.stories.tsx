import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

import * as RadioGroup from ".";

const options = [
  { label: "Public", value: "public" },
  { label: "Members", value: "members" },
  { label: "Private", value: "private" },
];

const meta = {
  title: "Components/Forms/Radio Group",
  component: RadioGroup.Root,
  parameters: {
    docs: {
      description: {
        component:
          "Presents a small set of mutually exclusive choices that should remain visible for comparison. Use Select when space is constrained or the list is longer, Checkbox for independent booleans, and Switch for an immediate on/off setting.",
      },
    },
  },
  args: {
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
    <LStack gap="6">
      {(["sm", "md", "lg"] as const).map((size) => (
        <LStack key={size} gap="2">
          <Text as="span" variant="metadata">
            {size}
          </Text>
          <Example size={size} defaultValue="public" />
        </LStack>
      ))}
      <LStack gap="2">
        <Text as="span" variant="metadata">
          Horizontal
        </Text>
        <Example size="md" orientation="horizontal" defaultValue="members" />
      </LStack>
    </LStack>
  ),
};

export const Disabled: Story = {
  render: () => (
    <RadioGroup.Root defaultValue="public">
      <RadioGroup.Item value="public">
        <RadioGroup.ItemControl />
        <RadioGroup.ItemText>Available option</RadioGroup.ItemText>
        <RadioGroup.ItemHiddenInput />
      </RadioGroup.Item>
      <RadioGroup.Item value="phone" disabled>
        <RadioGroup.ItemControl />
        <RadioGroup.ItemText>Unavailable option</RadioGroup.ItemText>
        <RadioGroup.ItemHiddenInput />
      </RadioGroup.Item>
    </RadioGroup.Root>
  ),
};
