import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, LStack, styled } from "@/styled-system/jsx";

import {
  TrailActionRunStatusBadge,
  TrailRunStatusBadge,
  TrailStatusBadge,
} from "./TrailStatusBadge";

const meta = {
  title: "Screens/Trails/Status badges",
  component: TrailStatusBadge,
  parameters: {
    docs: {
      description: {
        component:
          "Trail lifecycle, occurrence, and action statuses use explicit labels and semantic colour roles. Successful work is green, active work is informational blue, blocked work is warning amber, failures are danger red, and terminal neutral states remain gray.",
      },
    },
  },
} satisfies Meta<typeof TrailStatusBadge>;

export default meta;

type Story = StoryObj<typeof meta>;

export const AllStatuses: Story = {
  args: { status: "active" },
  render: () => (
    <LStack gap="4" alignItems="start">
      <BadgeGroup label="Trail">
        <TrailStatusBadge status="active" />
        <TrailStatusBadge status="paused" />
        <TrailStatusBadge status="finished" />
        <TrailStatusBadge status="archived" />
      </BadgeGroup>

      <BadgeGroup label="Run">
        <TrailRunStatusBadge status="queued" />
        <TrailRunStatusBadge status="running" />
        <TrailRunStatusBadge status="completed" />
        <TrailRunStatusBadge status="attention_required" />
        <TrailRunStatusBadge status="cancelled" />
        <TrailRunStatusBadge status="skipped" />
      </BadgeGroup>

      <BadgeGroup label="Action">
        <TrailActionRunStatusBadge status="queued" />
        <TrailActionRunStatusBadge status="running" />
        <TrailActionRunStatusBadge status="completed" />
        <TrailActionRunStatusBadge status="blocked" />
        <TrailActionRunStatusBadge status="failed" />
        <TrailActionRunStatusBadge status="cancelled" />
      </BadgeGroup>
    </LStack>
  ),
};

function BadgeGroup({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <LStack gap="1" alignItems="start">
      <styled.h2 fontSize="sm" fontWeight="semibold">
        {label}
      </styled.h2>
      <HStack gap="2" flexWrap="wrap">
        {children}
      </HStack>
    </LStack>
  );
}
