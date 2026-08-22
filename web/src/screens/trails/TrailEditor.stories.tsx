import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { TrailMutableProps } from "@/api/openapi-schema";
import { BackAction } from "@/components/site/Action/Back";
import { PageHeader } from "@/components/ui/page-header";
import { Box, LStack } from "@/styled-system/jsx";

import { TrailEditorForm } from "./TrailEditor";

const robots = [
  { id: "robot_moderator", name: "Community moderation assistant" },
  { id: "robot_curator", name: "Library directory curator" },
  { id: "robot_host", name: "Monthly photo thread host" },
];

const editValue = {
  name: "Scheduled moderation check",
  description: "Review new reports before the morning moderation shift.",
  status: "paused",
  trigger: {
    type: "schedule",
    schedule: {
      start: "2026-08-25T09:00:00",
      timezone: "Europe/London",
      rule: {
        frequency: "weekly",
        interval: 1,
        by_weekday: ["monday", "wednesday", "friday"],
      },
    },
  },
  actions: [
    {
      type: "robot_run",
      robot_ref: "robot_moderator",
      instruction:
        "Review unresolved reports, summarize any patterns, and finish with a structured result. Block if a moderator must make a policy decision.",
    },
  ],
} satisfies TrailMutableProps;

const occurrences = [
  "2026-08-24T08:00:00Z",
  "2026-08-26T08:00:00Z",
  "2026-08-28T08:00:00Z",
  "2026-08-31T08:00:00Z",
  "2026-09-02T08:00:00Z",
];

const meta = {
  title: "Screens/Trails/Editor",
  component: TrailEditorForm,
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "The Trail form is an administrative full-page workflow. The creator is always the signed-in account and is not editable. Name and Description open the form without a redundant section heading. Cadence uses Select; timezone and Robot choices use searchable Comboboxes; intervals use Number Input; dates use Date Picker; local time uses a time Input; weekdays use a content-sized Toggle Group. Timezone choices come from the browser's IANA list and the backend remains authoritative. The sole Create or Save action is right-aligned; page-level back navigation replaces a duplicate Cancel action.",
      },
    },
  },
  decorators: [
    (Story) => (
      <Box marginX="auto" maxW="5xl" padding={{ base: "4", md: "8" }}>
        <Story />
      </Box>
    ),
  ],
  args: {
    robots,
    onPreview: async () => occurrences,
    onSubmit: async () => undefined,
  },
} satisfies Meta<typeof TrailEditorForm>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Create: Story = {};

export const FullPageCreate: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "The complete creation composition uses a title-only Page Header, the back action as the exit path, contextual guidance beside fields, and one right-aligned confirmation action.",
      },
    },
  },
  render: (args) => (
    <LStack gap="4" alignItems="stretch">
      <PageHeader
        title="New Trail"
        back={<BackAction fallbackHref="/robots/trails" />}
      />
      <TrailEditorForm {...args} />
    </LStack>
  ),
};

export const EditPausedTrail: Story = {
  args: {
    initialValue: editValue,
  },
};

export const RobotsUnavailable: Story = {
  args: {
    robots: [],
    robotsError: "Could not load available Robots.",
  },
};
