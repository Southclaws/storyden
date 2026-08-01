import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { Button } from "../button";

import { Group } from ".";

const meta = {
  title: "UI/Group",
  component: Group,
  argTypes: {
    orientation: {
      control: "select",
      options: ["horizontal", "vertical"],
    },
    attached: {
      control: "boolean",
    },
    grow: {
      control: "boolean",
    },
  },
} satisfies Meta<typeof Group>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {
  args: {
    orientation: "horizontal",
    attached: false,
    grow: false,
  },
  render: (args) => (
    <Group {...args}>
      <Button variant="outline">Preview</Button>
      <Button variant="outline">Edit</Button>
      <Button variant="outline">Publish</Button>
    </Group>
  ),
};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      <Group orientation="horizontal">
        <Button variant="outline">One</Button>
        <Button variant="outline">Two</Button>
      </Group>
      <Group orientation="horizontal" attached>
        <Button variant="outline">One</Button>
        <Button variant="outline">Two</Button>
      </Group>
      <Group orientation="vertical" attached>
        <Button variant="outline">One</Button>
        <Button variant="outline">Two</Button>
      </Group>
    </LStack>
  ),
};
