import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { WarningIcon } from "../icons/Warning";

import * as Alert from ".";

const meta = {
  title: "Components/Feedback/Alert",
  component: Alert.Root,
  parameters: {
    docs: {
      description: {
        component:
          "Persistent contextual information that affects the current task. Alert uses Inset structure and adds status colour only to communicate meaning. Use it for risks, requirements, and important consequences; ordinary field guidance belongs in helper text, and transient operation feedback belongs in a toast.",
      },
    },
  },
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
