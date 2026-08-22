import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { LayoutGridIcon } from "../icons/LayoutGrid";
import { LayoutListIcon } from "../icons/LayoutList";
import { LayoutTableIcon } from "../icons/LayoutTable";

import * as ToggleGroup from ".";

const meta = {
  title: "Components/Forms/Toggle Group",
  component: ToggleGroup.Root,
  parameters: {
    docs: {
      description: {
        component:
          "Presents a compact set of closely related choices. The group sizes to its choices even inside a stretching form layout; do not distribute short labels across the available row. Use Select when the option set is longer or does not benefit from simultaneous visibility.",
      },
    },
  },
  args: {
    defaultValue: ["list"],
    size: "md",
    variant: "outline",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["outline", "ghost"],
    },
  },
} satisfies Meta<typeof ToggleGroup.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: ToggleGroup.RootProps) {
  return (
    <ToggleGroup.Root {...props}>
      <ToggleGroup.Item value="list" aria-label="List">
        <LayoutListIcon />
      </ToggleGroup.Item>
      <ToggleGroup.Item value="grid" aria-label="Grid">
        <LayoutGridIcon />
      </ToggleGroup.Item>
      <ToggleGroup.Item value="table" aria-label="Table">
        <LayoutTableIcon />
      </ToggleGroup.Item>
    </ToggleGroup.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      {(["xs", "sm", "md", "lg"] as const).map((size) => (
        <Example
          key={size}
          size={size}
          variant="outline"
          defaultValue={["list"]}
        />
      ))}
      <Example size="md" variant="ghost" defaultValue={["grid"]} />
    </LStack>
  ),
};

export const ContentSizedTextChoices: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "This intentionally sits in a full-width stretching stack. The Toggle Group remains content-sized instead of filling the row.",
      },
    },
  },
  render: () => (
    <LStack gap="2" alignItems="stretch" width="full">
      <strong>Weekdays</strong>
      <ToggleGroup.Root
        defaultValue={["sunday"]}
        multiple
        size="sm"
        variant="outline"
      >
        {(
          [
            ["monday", "mon"],
            ["tuesday", "tue"],
            ["wednesday", "wed"],
            ["thursday", "thu"],
            ["friday", "fri"],
            ["saturday", "sat"],
            ["sunday", "sun"],
          ] as const
        ).map(([value, label]) => (
          <ToggleGroup.Item key={value} value={value} aria-label={value}>
            {label}
          </ToggleGroup.Item>
        ))}
      </ToggleGroup.Root>
    </LStack>
  ),
};
