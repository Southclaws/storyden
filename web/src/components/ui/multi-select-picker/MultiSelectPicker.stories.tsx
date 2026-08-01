import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { MultiSelectPicker, type MultiSelectPickerItem } from ".";

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

const meta = {
  title: "UI/Multi Select Picker",
  component: MultiSelectPicker,
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
      options: ["xs", "sm", "md", "lg", "xl", "2xl"],
    },
  },
} satisfies Meta<typeof MultiSelectPicker>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: (args) => {
    const [value, setValue] = useState<MultiSelectPickerItem[]>(selectedItems);
    const [results, setResults] =
      useState<MultiSelectPickerItem[]>(queryResults);

    return (
      <MultiSelectPicker
        {...args}
        value={value}
        queryResults={results}
        onChange={async (next) => setValue(next)}
        onQuery={(query) =>
          setResults(
            queryResults.filter((item) =>
              item.label.toLowerCase().includes(query.toLowerCase()),
            ),
          )
        }
      />
    );
  },
};
