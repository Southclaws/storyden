import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Button } from "../button";
import { CheckIcon } from "../icons/Check";
import { DeleteIcon } from "../icons/Delete";
import { EditIcon } from "../icons/Edit";

import * as Menu from ".";

const meta = {
  title: "UI/Menu",
  component: Menu.Root,
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof Menu.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: Menu.RootProps) {
  return (
    <Menu.Root {...props}>
      <Menu.Trigger asChild>
        <Button variant="outline" size={props.size}>
          Open menu
        </Button>
      </Menu.Trigger>
      <Menu.Positioner>
        <Menu.Content>
          <Menu.ItemGroup>
            <Menu.ItemGroupLabel>Actions</Menu.ItemGroupLabel>
            <Menu.Item value="edit">
              <EditIcon />
              <Menu.ItemText>Edit</Menu.ItemText>
            </Menu.Item>
            <Menu.Item value="publish">
              <CheckIcon />
              <Menu.ItemText>Publish</Menu.ItemText>
            </Menu.Item>
            <Menu.Item value="delete">
              <DeleteIcon />
              <Menu.ItemText>Delete</Menu.ItemText>
            </Menu.Item>
          </Menu.ItemGroup>
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  );
}

export const Playground: Story = {
  args: {
    size: "xs",
  },
  render: (args) => <Example {...args} />,
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      {(["xs", "sm", "md", "lg"] as const).map((size) => (
        <Example key={size} size={size} />
      ))}
    </LStack>
  ),
};
