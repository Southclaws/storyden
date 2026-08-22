import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Trail, TrailRun } from "@/api/openapi-schema";
import { Box, LStack } from "@/styled-system/jsx";

import { TrailCard } from "./TrailListScreen";

const now = new Date();

const baseTrail = {
  id: "trail_moderation",
  createdAt: minutesAgo(240),
  updatedAt: minutesAgo(12),
  name: "Community moderation",
  description: "Review new reports and flag decisions that need a moderator.",
  created_by: {
    id: "account_admin",
    joined: minutesAgo(1440),
    handle: "operator",
    name: "Community operator",
    roles: [],
  },
  status: "active",
  trigger: {
    type: "schedule",
    schedule: {
      start: "2026-08-22T09:00:00",
      timezone: "Europe/London",
      rule: {
        frequency: "hourly",
        interval: 1,
      },
    },
  },
  actions: [action("action_review")],
  next_occurrence_at: minutesFromNow(48),
  last_occurrence_at: minutesAgo(12),
} satisfies Trail;

const completedRun = run("completed", 12);

const meta = {
  title: "Screens/Trails/List card",
  component: TrailCard,
  parameters: {
    docs: {
      description: {
        component:
          "Trail cards are compact CardBox compositions and operational health summaries. Identity occupies the upper-left, current execution or lifecycle state occupies the upper-right, generic action count occupies the lower-left, and latest activity occupies the lower-right. Successful activity stays quiet. Running and attention-required states receive stronger semantic treatment. Cadence is not used as a fallback description because future Trails can use non-schedule triggers. The whole surface navigates to the Trail and contains no action menu.",
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
    trail: baseTrail,
    latest: completedRun,
  },
} satisfies Meta<typeof TrailCard>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Completed: Story = {};

export const UnbrokenContent: Story = {
  args: {
    trail: {
      ...baseTrail,
      id: "trail_unbroken_content",
      name: "W".repeat(120),
      description: "W".repeat(120),
    },
  },
  parameters: {
    docs: {
      description: {
        story:
          "An unbroken 120-character name and description stay within the identity column without obscuring operational state or activity metadata.",
      },
    },
  },
};

export const OperationalStates: Story = {
  render: () => (
    <LStack gap="3" alignItems="stretch">
      <TrailCard trail={baseTrail} latest={completedRun} />
      <TrailCard
        trail={{
          ...baseTrail,
          id: "trail_attention",
          name: "New member safety review",
          description: "Check recent sign-ups for moderation concerns.",
          actions: [action("action_review"), action("action_notify")],
        }}
        latest={run("attention_required", 8)}
      />
      <TrailCard
        trail={{
          ...baseTrail,
          id: "trail_running",
          name: "Library directory refresh",
          description: "Review new Library pages and update the directory.",
          actions: [
            action("action_review"),
            action("action_update"),
            action("action_report"),
          ],
        }}
        latest={run("running", 2)}
      />
      <TrailCard
        trail={{
          ...baseTrail,
          id: "trail_paused",
          name: "Monthly community prompt",
          description: "Invite members to share their latest photos.",
          status: "paused",
          next_occurrence_at: undefined,
        }}
        latest={run("completed", 180)}
      />
      <TrailCard
        trail={{
          ...baseTrail,
          id: "trail_new",
          name: "Weekly digest",
          description: "Summarise the discussions members followed this week.",
          last_occurrence_at: undefined,
        }}
      />
    </LStack>
  ),
};

export const ActivityLoading: Story = {
  args: {
    latest: undefined,
    activityState: "loading",
  },
};

export const ActivityUnavailable: Story = {
  args: {
    latest: undefined,
    activityState: "failed",
  },
};

function action(id: string): Trail["actions"][number] {
  return {
    id,
    createdAt: minutesAgo(240),
    updatedAt: minutesAgo(240),
    action: {
      type: "robot_run",
      robot_ref: "robot_moderator",
      instruction: "Review recent activity and report anything unusual.",
    },
  };
}

function run(status: TrailRun["status"], minutes: number): TrailRun {
  const timestamp = minutesAgo(minutes);

  return {
    id: `run_${status}`,
    createdAt: timestamp,
    updatedAt: timestamp,
    trail_id: baseTrail.id,
    trigger: {
      trail_id: baseTrail.id,
      trail_run_id: `run_${status}`,
      kind: "scheduled",
      trigger: baseTrail.trigger,
      scheduled_for: timestamp,
      observed_at: timestamp,
    },
    status,
    scheduled_for: timestamp,
    actions: [],
    finished_at:
      status === "running" || status === "queued" ? undefined : timestamp,
  };
}

function minutesAgo(minutes: number): string {
  return new Date(now.getTime() - minutes * 60_000).toISOString();
}

function minutesFromNow(minutes: number): string {
  return new Date(now.getTime() + minutes * 60_000).toISOString();
}
