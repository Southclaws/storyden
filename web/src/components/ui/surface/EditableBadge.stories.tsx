import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";
import { fn } from "storybook/test";

import { Grid, HStack, LStack } from "@/styled-system/jsx";

import { Text } from "../text";

import { EditableBadge } from ".";

const meta = {
  title: "Components/Data Display/Surface/Editable Badge",
  component: EditableBadge,
  parameters: {
    nextjs: { appDirectory: false },
    docs: {
      description: {
        component:
          "Inline editing for short Card metadata. View and edit modes share the canonical Badge shell and stay on one line, truncating to the parent layout. The native input scrolls horizontally while editing. Hover and focus reveal editability without changing the pill's dimensions. Enter finishes editing; IME composition remains uninterrupted. Empty fields retain a named placeholder. Use Input for ordinary form fields and Badge for metadata that cannot be edited.",
      },
    },
  },
  args: {
    editing: true,
    value: "Commercial",
    inputLabel: "License",
    placeholder: "License",
    onValueChange: fn(),
    size: "sm",
    colorPalette: "gray",
    h: "auto",
    py: "0.5",
    maxW: "full",
  },
} satisfies Meta<typeof EditableBadge>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Comparison: Story = {
  render: function Comparison(args) {
    const [value, setValue] = useState(args.value);
    return (
      <HStack gap="6" alignItems="start" maxW="sm">
        <LStack gap="2" minW="0">
          <Text variant="metadata">View</Text>
          <EditableBadge
            {...args}
            editing={false}
            value={value}
            onValueChange={setValue}
          />
        </LStack>
        <LStack gap="2" minW="0">
          <Text variant="metadata">Edit</Text>
          <EditableBadge {...args} value={value} onValueChange={setValue} />
        </LStack>
      </HStack>
    );
  },
};

export const Empty: Story = { ...Comparison, args: { value: "" } };
export const LongValue: Story = {
  ...Comparison,
  args: { value: "A long classification that stays on one line" },
};

export const CompetingValues: Story = {
  args: { value: `Monospace ${"W".repeat(200)}` },
  render: function CompetingValues(args) {
    const [style, setStyle] = useState(args.value);
    const [use, setUse] = useState("Code");
    return (
      <Grid
        w="full"
        maxW="sm"
        minW="0"
        gridTemplateColumns="minmax(0, max-content) minmax(0, max-content)"
        justifyContent="space-between"
        gap="2"
      >
        <HStack minW="0" flexWrap="wrap" gap="1">
          <EditableBadge
            {...args}
            value={style}
            inputLabel="Style"
            onValueChange={setStyle}
          />
        </HStack>
        <HStack minW="0" flexWrap="wrap" justify="end" gap="1">
          <EditableBadge
            {...args}
            value={use}
            inputLabel="Use"
            onValueChange={setUse}
          />
        </HStack>
      </Grid>
    );
  },
};
