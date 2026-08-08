import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box, LStack } from "@/styled-system/jsx";

import { Button } from "../button";
import { CheckIcon } from "../icons/Check";
import { DeleteIcon } from "../icons/Delete";
import { EditIcon } from "../icons/Edit";

import * as Menu from ".";
import { Menu as ClosedMenu } from "./Menu";

const meta = {
  title: "Components/Overlays/Menu",
  component: ClosedMenu,
  args: {
    trigger: <Button variant="outline">Open menu</Button>,
    children: <Menu.Item value="example">Example</Menu.Item>,
  },
  argTypes: {
    size: {
      control: "select",
      options: ["xs", "sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof ClosedMenu>;

export default meta;

type Story = StoryObj<typeof meta>;

function Example(props: Menu.RootProps) {
  const triggerSize =
    props.size === "md" || props.size === "lg" ? props.size : "sm";

  return (
    <ClosedMenu
      {...props}
      trigger={
        <Button variant="outline" size={triggerSize}>
          Open menu
        </Button>
      }
    >
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
    </ClosedMenu>
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

export const ViewportEdge: Story = {
  render: () => (
    <Box display="flex" justifyContent="end" minH="sm" width="full">
      <Example positioning={{ placement: "bottom-end" }} />
    </Box>
  ),
};

export const Compound: Story = {
  render: () => (
    <Menu.Root>
      <Menu.Trigger asChild>
        <Button variant="outline">Custom menu</Button>
      </Menu.Trigger>
      <Menu.Positioner>
        <Menu.Content>
          <Menu.Item value="custom">Custom composition</Menu.Item>
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  ),
};
