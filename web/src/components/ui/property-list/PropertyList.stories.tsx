import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { Input } from "@/components/ui/input";

import { PropertyList } from ".";

const meta = {
  title: "Components/Data Display/Property List",
  component: PropertyList,
  parameters: {
    nextjs: { appDirectory: false },
    docs: {
      description: {
        component:
          "Compact metadata for one object, expressed as a description list with shared label and value columns. Use in cards and object summaries. Labels size to their content up to 40% of the list; values use the remaining width, align right, and clamp to two lines. Single-line editors retain the full value and scroll horizontally. Label and value slots accept controls without losing dt/dd semantics. Use Table for comparing records across rows, Badge for a short classification, and ordinary form fields for a standalone form.",
      },
    },
  },
  args: {
    maxW: "sm",
    items: [
      { id: "style", label: "Style", value: "Monospace" },
      { id: "license", label: "License", value: "Open source" },
      { id: "use", label: "Use", value: "Code" },
    ],
  },
} satisfies Meta<typeof PropertyList>;
export default meta;
type Story = StoryObj<typeof meta>;

export const ReadOnly: Story = {};
export const Empty: Story = { args: { items: [] } };
export const LongValues: Story = {
  args: {
    items: [
      { id: "style", label: "Style", value: `Monospace ${"W".repeat(200)}` },
      {
        id: "license",
        label: "License",
        value: `Open source ${"W".repeat(200)}`,
      },
      {
        id: "use",
        label: "A longer property name with context",
        value: "Code",
      },
    ],
  },
};

export const Editable: Story = {
  ...LongValues,
  render: function Editable(args) {
    const [items, setItems] = useState(args.items);
    return (
      <PropertyList
        {...args}
        items={items.map((item) => ({
          ...item,
          value: (
            <Input
              aria-label={String(item.label)}
              variant="ghost"
              size="sm"
              minW="0"
              textAlign="inherit"
              value={String(item.value)}
              _focus={{ outline: "none" }}
              _focusVisible={{ outline: "none" }}
              onChange={(event) =>
                setItems((current) =>
                  current.map((entry) =>
                    entry.id === item.id
                      ? { ...entry, value: event.target.value }
                      : entry,
                  ),
                )
              }
            />
          ),
        }))}
      />
    );
  },
};
