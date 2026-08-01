import { createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { CheckIcon } from "../icons/Check";
import { SelectIcon } from "../icons/Select";
import { Input } from "../input";

import * as Combobox from ".";

const items = [
  { label: "General", value: "general" },
  { label: "Announcements", value: "announcements" },
  { label: "Support", value: "support" },
  { label: "Showcase", value: "showcase" },
];

const collection = createListCollection({ items });

type ComboboxStoryArgs = {
  size?: "xs" | "sm" | "md" | "lg";
};

const meta = {
  title: "UI/Combobox",
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg"],
    },
  },
} satisfies Meta<ComboboxStoryArgs>;

export default meta;

type Story = StoryObj<ComboboxStoryArgs>;

function Example(props: ComboboxStoryArgs) {
  return (
    <Combobox.Root collection={collection} maxW="sm" size={props.size}>
      <Combobox.Label>Category</Combobox.Label>
      <Combobox.Control>
        <Combobox.Input placeholder="Select a category" asChild>
          <Input size={props.size} />
        </Combobox.Input>
        <Combobox.Trigger asChild>
          <button type="button" aria-label="Open">
            <SelectIcon />
          </button>
        </Combobox.Trigger>
      </Combobox.Control>
      <Combobox.Positioner>
        <Combobox.Content>
          <Combobox.List>
            <Combobox.ItemGroup>
              {items.map((item) => (
                <Combobox.Item key={item.value} item={item}>
                  <Combobox.ItemText>{item.label}</Combobox.ItemText>
                  <Combobox.ItemIndicator>
                    <CheckIcon />
                  </Combobox.ItemIndicator>
                </Combobox.Item>
              ))}
            </Combobox.ItemGroup>
          </Combobox.List>
        </Combobox.Content>
      </Combobox.Positioner>
    </Combobox.Root>
  );
}

export const Playground: Story = {
  args: {
    size: "md",
  },
  render: (args) => <Example {...args} />,
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4">
      {(["xs", "sm", "md", "lg"] as const).map((size) => (
        <Example key={size} size={size} />
      ))}
    </LStack>
  ),
};
