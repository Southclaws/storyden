import { createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { CheckIcon } from "../icons/Check";
import { SelectIcon } from "../icons/Select";

import * as Select from ".";

const items = [
  { label: "Discussion", value: "discussion" },
  { label: "Library", value: "library" },
  { label: "Collection", value: "collection" },
];

const collection = createListCollection({ items });

type SelectStoryArgs = {
  size?: "sm" | "md" | "lg";
  variant?: "outline" | "ghost";
};

const meta = {
  title: "Components/Forms/Select",
  args: {
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
      options: ["outline", "ghost"],
    },
  },
} satisfies Meta<SelectStoryArgs>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: SelectStoryArgs) {
  return (
    <Select.Root
      collection={collection}
      maxW="sm"
      size={props.size}
      variant={props.variant}
    >
      <Select.Label>Content type</Select.Label>
      <Select.Control>
        <Select.Trigger>
          <Select.ValueText placeholder="Select content type" />
          <SelectIcon />
        </Select.Trigger>
      </Select.Control>
      <Select.Positioner>
        <Select.Content>
          {items.map((item) => (
            <Select.Item key={item.value} item={item}>
              <Select.ItemText>{item.label}</Select.ItemText>
              <Select.ItemIndicator>
                <CheckIcon />
              </Select.ItemIndicator>
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Positioner>
      <Select.HiddenSelect />
    </Select.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Variants: Story = {
  render: () => (
    <LStack gap="4">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Example key={size} size={size} />
      ))}
      <Example variant="ghost" />
    </LStack>
  ),
};
