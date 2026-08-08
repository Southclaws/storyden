import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { WarningIcon } from "../icons/Warning";

import * as Alert from ".";

const meta = {
  title: "Components/Feedback/Alert",
  component: Alert.Root,
} satisfies Meta<typeof Alert.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Alert.Root maxW="lg">
      <Alert.Icon asChild>
        <WarningIcon />
      </Alert.Icon>
      <Alert.Content>
        <Alert.Title>Moderation queue is busy</Alert.Title>
        <Alert.Description>
          New reports are still being accepted while the current review batch is
          processed.
        </Alert.Description>
      </Alert.Content>
    </Alert.Root>
  ),
};
