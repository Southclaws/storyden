import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { TrailIcon } from "@/components/ui/icons/Trail";

import { EmptyState } from "./EmptyState";

const meta = {
  title: "Compositions/Feedback/Empty State",
  component: EmptyState,
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component:
          "Centres an empty state and explains the next useful action. Use plain product language and distinguish first use from no results, permissions, and failures. Operational and administrative screens must hide the generic contribution label; contribution language is only appropriate for community content surfaces.",
      },
    },
  },
} satisfies Meta<typeof EmptyState>;

export default meta;

type Story = StoryObj<typeof meta>;

export const OperationalFirstUse: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Operational empty states say what is missing and what the administrator can do next. They do not expose implementation terms or invite contribution.",
      },
    },
  },
  render: () => (
    <EmptyState icon={<TrailIcon />} hideContributionLabel minH="64">
      No Trails yet. Create one to run Robots automatically on a schedule.
    </EmptyState>
  ),
};

export const NoFilteredResults: Story = {
  render: () => (
    <EmptyState hideContributionLabel minH="64">
      No Trails match the current filters.
    </EmptyState>
  ),
};
