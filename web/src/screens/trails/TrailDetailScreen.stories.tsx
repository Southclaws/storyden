import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Trail, TrailActionRun, TrailRun } from "@/api/openapi-schema";
import { Box } from "@/styled-system/jsx";

import { TrailDetailView } from "./TrailDetailScreen";

const now = new Date("2026-08-22T16:00:00.000Z");

const robots = [
  { id: "robot_curator", name: "Community directory curator" },
  { id: "robot_moderator", name: "Moderation assistant" },
];

const baseTrail = {
  id: "trail_hourly_check",
  createdAt: minutesAgo(240),
  updatedAt: minutesAgo(12),
  name: "Hourly community check",
  description:
    "Review recent activity and surface anything that needs an operator.",
  created_by: {
    id: "account_operator",
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
  actions: [
    binding(
      "binding_curator",
      "robot_curator",
      "Review recent Library activity and summarise notable changes.",
    ),
    binding(
      "binding_moderator",
      "robot_moderator",
      "Review moderation activity and report anything that needs attention.",
    ),
  ],
  next_occurrence_at: minutesFromNow(48),
  last_occurrence_at: minutesAgo(12),
} satisfies Trail;

const completedRun = run({
  id: "run_completed",
  kind: "scheduled",
  status: "completed",
  minutes: 12,
  actions: [
    actionRun({
      id: "action_completed",
      binding: baseTrail.actions[0]!,
      status: "completed",
      summary:
        "The Library directory is current. Three new pages were reviewed and no action is needed.",
      sessionID: "session_completed",
    }),
  ],
});

const meta = {
  title: "Screens/Trails/Detail",
  component: TrailDetailView,
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "Trail detail keeps Run now visible as the primary mobile action and moves Edit, Pause or Resume, and Archive into an overflow menu. The menu trigger and Run now form one attached action group. Desktop retains the full action group. The compact mobile controls share the title row without crushing identity. The description and lifecycle status belong to the Trail overview instead of competing with identity. Creator and action count form one metadata group, with lifecycle status isolated in the opposite corner. A curated schedule overview groups trigger identity and temporal state by meaning instead of presenting them as equal metrics. The schedule icon sits inline with its label so cadence, timezone, and occurrence details share one left edge at every width. Run history leads with relative time, then gives the trigger type and precise occurrence time as supporting context. Aggregate run status occupies the opposite corner. A single action omits its duplicate status badge; runs with multiple independent actions retain each action status. User-facing cards never expose resource identifiers. Robot sessions use a visible LinkButton. The view is resilient to long unbroken identity, instruction, result, and error content.",
      },
    },
  },
  decorators: [
    (Story) => (
      <Box maxW="7xl" marginX="auto" padding={{ base: "4", md: "8" }}>
        <Story />
      </Box>
    ),
  ],
  args: {
    trail: baseTrail,
    runs: [completedRun],
    robots,
    onRunNow: async () => undefined,
    onChangeState: async () => undefined,
    onCancelAction: async () => undefined,
  },
} satisfies Meta<typeof TrailDetailView>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Completed: Story = {};

export const EventTriggered: Story = {
  args: {
    trail: {
      ...baseTrail,
      id: "trail_event_check",
      name: "Published content check",
      trigger: {
        type: "event",
        events: ["EventThreadPublished", "EventNodePublished"],
      },
      next_occurrence_at: undefined,
    },
  },
};

export const NeedsAttention: Story = {
  args: {
    runs: [
      run({
        id: "run_attention",
        kind: "manual",
        status: "attention_required",
        minutes: 4,
        actions: [
          actionRun({
            id: "action_success",
            binding: baseTrail.actions[0]!,
            status: "completed",
            summary: "The Library review completed without changes.",
            sessionID: "session_success",
          }),
          actionRun({
            id: "action_blocked",
            binding: baseTrail.actions[1]!,
            status: "blocked",
            summary: "The moderation review needs an operator decision.",
            attention: {
              reason: "needs_approval",
              message: "Approve the suggested moderation action to continue.",
            },
            sessionID: "session_blocked",
          }),
        ],
      }),
      completedRun,
    ],
  },
};

export const Running: Story = {
  args: {
    runs: [
      run({
        id: "run_running",
        kind: "manual",
        status: "running",
        minutes: 1,
        actions: [
          actionRun({
            id: "action_running",
            binding: baseTrail.actions[0]!,
            status: "running",
            sessionID: "session_running",
          }),
          actionRun({
            id: "action_queued",
            binding: baseTrail.actions[1]!,
            status: "queued",
          }),
        ],
      }),
    ],
  },
};

export const EmptyHistory: Story = {
  args: {
    runs: [],
  },
};

export const UnavailableTriggerDetails: Story = {
  args: {
    runs: [
      {
        ...completedRun,
        id: "run_unavailable_trigger",
        trigger: undefined,
      },
    ],
  },
};

export const LongContent: Story = {
  args: {
    trail: {
      ...baseTrail,
      id: "trail_long_content",
      name: "W".repeat(120),
      description: "W".repeat(480),
      actions: [
        binding(
          "binding_long",
          "robot_curator",
          `Check this uninterrupted instruction: ${"W".repeat(240)}`,
        ),
      ],
    },
    runs: [
      run({
        id: "run_long_content",
        kind: "scheduled",
        status: "completed",
        minutes: 12,
        actions: [
          actionRun({
            id: "action_long_content",
            binding: binding(
              "binding_long_snapshot",
              "robot_curator",
              `Check this uninterrupted instruction: ${"W".repeat(240)}`,
            ),
            status: "completed",
            summary: `The Robot returned an uninterrupted result: ${"W".repeat(320)}`,
            sessionID: "session_long_content",
          }),
        ],
      }),
    ],
  },
};

function binding(
  id: string,
  robotRef: string,
  instruction: string,
): Trail["actions"][number] {
  return {
    id,
    createdAt: minutesAgo(240),
    updatedAt: minutesAgo(240),
    action: {
      type: "robot_run",
      robot_ref: robotRef,
      instruction,
    },
  };
}

function actionRun({
  id,
  binding: action,
  status,
  summary,
  attention,
  error,
  sessionID,
}: {
  id: string;
  binding: Trail["actions"][number];
  status: TrailActionRun["status"];
  summary?: string;
  attention?: { reason: string; message: string };
  error?: string;
  sessionID?: string;
}): TrailActionRun {
  const timestamp = minutesAgo(4);

  return {
    id,
    createdAt: timestamp,
    updatedAt: timestamp,
    action,
    status,
    error,
    started_at: timestamp,
    finished_at:
      status === "queued" || status === "running" ? undefined : timestamp,
    target: sessionID
      ? {
          type: "robot_run",
          robot_session_id: sessionID,
          output: summary
            ? {
                status:
                  status === "blocked"
                    ? "blocked"
                    : status === "failed"
                      ? "failed"
                      : "completed",
                summary,
                attention,
              }
            : undefined,
        }
      : undefined,
  };
}

function run({
  id,
  kind,
  status,
  minutes,
  actions,
}: {
  id: string;
  kind: NonNullable<TrailRun["trigger"]>["kind"];
  status: TrailRun["status"];
  minutes: number;
  actions: TrailActionRun[];
}): TrailRun {
  const timestamp = minutesAgo(minutes);

  return {
    id,
    createdAt: timestamp,
    updatedAt: timestamp,
    trail_id: baseTrail.id,
    trigger: {
      trail_id: baseTrail.id,
      trail_run_id: id,
      kind,
      trigger: baseTrail.trigger,
      scheduled_for: kind === "scheduled" ? timestamp : undefined,
      observed_at: timestamp,
    },
    status,
    scheduled_for: kind === "scheduled" ? timestamp : undefined,
    actions,
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
