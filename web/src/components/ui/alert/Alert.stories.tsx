import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { WarningIcon } from "../icons/Warning";

import * as Alert from ".";

const meta = {
  title: "Components/Feedback/Alert",
  component: Alert.Root,
  parameters: {
    docs: {
      description: {
        component:
          "Persistent contextual information that affects the current task. Alert uses Inset structure and adds status colour only to communicate meaning. Use it instead of assembling a locally coloured Box or Stack for risks, requirements, and important consequences. Ordinary field guidance belongs in helper text, and transient operation feedback belongs in a toast.",
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

export const SemanticTones: Story = {
  render: () => (
    <LStack gap="3" maxW="lg">
      {(["info", "success", "warning", "danger"] as const).map((tone) => (
        <Alert.Root key={tone} tone={tone}>
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>{tone} alert</Alert.Title>
            <Alert.Description>
              The surface, border, and content share one semantic status tone.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>
      ))}
    </LStack>
  ),
};
