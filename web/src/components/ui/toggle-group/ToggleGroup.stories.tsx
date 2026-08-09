import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { LayoutGridIcon } from "../icons/LayoutGrid";
import { LayoutListIcon } from "../icons/LayoutList";
import { LayoutTableIcon } from "../icons/LayoutTable";

import * as ToggleGroup from ".";

const meta = {
  title: "Components/Forms/Toggle Group",
  component: ToggleGroup.Root,
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
