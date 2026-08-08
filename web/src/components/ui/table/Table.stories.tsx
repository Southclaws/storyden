import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import * as Table from ".";

const rows = [
  ["Storybook", "In progress", "Design system"],
  ["Navigation shell", "Done", "Layout"],
  ["Theme hooks", "Next", "Customization"],
];

const meta = {
  title: "Components/Data Display/Table",
  component: Table.Root,
  args: {
    size: "md",
    variant: "plain",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md"],
    },
    variant: {
      control: "select",
      options: ["plain", "dense"],
    },
  },
} satisfies Meta<typeof Table.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: Table.RootProps) {
  return (
    <Table.Root {...props}>
      <Table.Caption>Component work queue</Table.Caption>
      <Table.Head>
        <Table.Row>
          <Table.Header>Task</Table.Header>
          <Table.Header>Status</Table.Header>
          <Table.Header>Area</Table.Header>
        </Table.Row>
      </Table.Head>
      <Table.Body>
        {rows.map(([task, status, area]) => (
          <Table.Row key={task}>
            <Table.Cell>{task}</Table.Cell>
            <Table.Cell>{status}</Table.Cell>
            <Table.Cell>{area}</Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Variants: Story = {
  render: () => (
    <LStack gap="6">
      <Example size="sm" variant="plain" />
      <Example size="md" variant="dense" />
    </LStack>
  ),
};
