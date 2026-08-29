import { createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { FormControl } from "@/components/ui/form-control";
import { FormLabel } from "@/components/ui/form-label";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import {
  MultiSelectPicker,
  type MultiSelectPickerItem,
  type MultiSelectPickerProps,
} from ".";

const selectedItems = [
  { label: "Design", value: "design", colour: "#6d5dfc" },
  { label: "Frontend", value: "frontend", colour: "#27b981" },
];

const queryResults = [
  { label: "Accessibility", value: "accessibility" },
  { label: "Moderation", value: "moderation" },
  { label: "Performance", value: "performance" },
  { label: "Storybook", value: "storybook" },
];

const categories = createListCollection({
  items: [
    { label: "Discussion", value: "discussion" },
    { label: "Question", value: "question" },
  ],
});

const meta = {
  title: "Components/Forms/Multi Select Picker",
  component: MultiSelectPicker,
  parameters: {
    docs: {
      description: {
        component:
          "A searchable multi-value combobox for choosing several items from a dynamic collection. Use it when selection and search belong in one compact field. Use Combobox for one value, Select for a short known list, and Checkbox Group when every choice should remain visible. Do not use it for commands, navigation, or a single free-form value.",
      },
    },
  },
  args: {
    value: selectedItems,
    queryResults,
    inputPlaceholder: "Select tags...",
    allowNewValues: true,
    autoColour: true,
    size: "sm",
    onChange: async () => undefined,
    onQuery: () => undefined,
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof MultiSelectPicker>;

export default meta;

type Story = StoryObj<typeof meta>;

type ControlledPickerProps = Omit<
  MultiSelectPickerProps,
  "onChange" | "onQuery" | "queryResults" | "value"
> & {
  initialValue?: MultiSelectPickerItem[];
  items?: MultiSelectPickerItem[];
};

function ControlledPicker({
  initialValue = selectedItems,
  items = queryResults,
  ...props
}: ControlledPickerProps) {
  const [value, setValue] = useState<MultiSelectPickerItem[]>(initialValue);
  const [results, setResults] = useState(items);

  return (
    <MultiSelectPicker
      {...props}
      value={value}
      queryResults={results}
      onChange={async (next) => setValue(next)}
      onQuery={(query) =>
        setResults(
          items.filter((item) =>
            item.label.toLowerCase().includes(query.toLowerCase()),
          ),
        )
      }
    />
  );
}

export const Default: Story = {
  render: (args) => (
    <LStack width="full" maxW="lg">
      <ControlledPicker
        {...args}
        initialValue={args.value}
        items={args.queryResults}
      />
    </LStack>
  ),
};

export const Sizes: Story = {
  render: () => (
    <LStack width="full" maxW="lg" gap="4">
      {(["sm", "md", "lg"] as const).map((size) => (
        <FormControl key={size}>
          <FormLabel>{size.toUpperCase()}</FormLabel>
          <ControlledPicker size={size} inputPlaceholder="Select tags..." />
        </FormControl>
      ))}
    </LStack>
  ),
};

export const AdjacentControlAlignment: Story = {
  render: () => (
    <HStack alignItems="center" width="full" maxW="lg">
      <Select.Root collection={categories} size="sm" width="48">
        <Select.Label srOnly>Category</Select.Label>
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder="Category" />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {categories.items.map((item) => (
              <Select.Item key={item.value} item={item}>
                <Select.ItemText>{item.label}</Select.ItemText>
              </Select.Item>
            ))}
          </Select.Content>
        </Select.Positioner>
        <Select.HiddenSelect />
      </Select.Root>

      <styled.div flex="1" minW="0">
        <ControlledPicker
          initialValue={[{ label: "tag", value: "tag" }]}
          inputPlaceholder="Add tags..."
        />
      </styled.div>
    </HStack>
  ),
};

export const NarrowWithManyValues: Story = {
  render: () => (
    <LStack width="64">
      <ControlledPicker
        initialValue={[...selectedItems, ...queryResults]}
        items={[]}
        inputPlaceholder="Add another tag..."
      />
    </LStack>
  ),
};

export const WideField: Story = {
  render: () => (
    <LStack width="full" maxW="5xl">
      <ControlledPicker defaultOpen inputPlaceholder="Select capabilities..." />
    </LStack>
  ),
};

export const NearViewportEdge: Story = {
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        story:
          "The popup flips above the field on its first open when the viewport edge leaves insufficient room below.",
      },
    },
  },
  render: () => (
    <styled.div alignItems="flex-end" display="flex" minH="screen" p="4">
      <LStack width="full" maxW="lg">
        <ControlledPicker inputPlaceholder="Select capabilities..." />
      </LStack>
    </styled.div>
  ),
};

export const Error: Story = {
  render: () => (
    <LStack width="full" maxW="lg">
      <ControlledPicker
        defaultOpen
        items={[]}
        queryError="Capabilities could not be loaded. Try again."
      />
    </LStack>
  ),
};

export const Disabled: Story = {
  render: () => (
    <LStack width="full" maxW="lg">
      <ControlledPicker disabled />
    </LStack>
  ),
};
