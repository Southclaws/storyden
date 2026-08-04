import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack } from "@/styled-system/jsx";

import { AddIcon } from "../icons/Add";
import { DeleteIcon } from "../icons/Delete";
import { EditIcon } from "../icons/Edit";

import { IconButton } from ".";

const meta = {
  title: "UI/Icon Button",
  component: IconButton,
  args: {
    "aria-label": "Create",
    children: <AddIcon />,
    size: "sm",
    variant: "subtle",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
    variant: {
      control: "select",
      options: ["solid", "outline", "ghost", "subtle", "plain"],
    },
  },
} satisfies Meta<typeof IconButton>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Variants: Story = {
  render: () => (
    <LStack gap="4" alignItems="start">
      <HStack gap="3" flexWrap="wrap">
        <IconButton aria-label="Add" variant="solid">
          <AddIcon />
        </IconButton>
        <IconButton aria-label="Edit" variant="subtle">
          <EditIcon />
        </IconButton>
        <IconButton aria-label="Delete" variant="outline">
          <DeleteIcon />
        </IconButton>
        <IconButton aria-label="Edit ghost" variant="ghost">
          <EditIcon />
        </IconButton>
        <IconButton aria-label="Edit plain" variant="plain">
          <EditIcon />
        </IconButton>
      </HStack>
      <HStack gap="3" flexWrap="wrap">
        {(["sm", "md", "lg"] as const).map((size) => (
          <IconButton key={size} size={size} aria-label={`Add ${size}`}>
            <AddIcon />
          </IconButton>
        ))}
      </HStack>
    </LStack>
  ),
};

export const States: Story = {
  render: () => (
    <HStack gap="3" flexWrap="wrap">
      <IconButton aria-label="Default">
        <AddIcon />
      </IconButton>
      <IconButton aria-label="Loading" loading>
        <AddIcon />
      </IconButton>
      <IconButton aria-label="Disabled" disabled>
        <AddIcon />
      </IconButton>
    </HStack>
  ),
};
