import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Badge } from "@/components/ui/badge";
import { HStack, LStack } from "@/styled-system/jsx";

import { StatusBadge, type StatusBadgeTone } from ".";

const meta = {
  title: "Components/Data Display/Status Badge",
  component: StatusBadge,
  parameters: {
    docs: {
      description: {
        component:
          "Communicates the state of an object or operation through a stable semantic tone. Use it instead of manually styling Badge with status token triples. Domain components remain responsible for mapping their statuses to a label and tone. Use plain Badge for categories, tags, and other labels that are not lifecycle or execution state.",
      },
    },
  },
  args: {
    children: "Completed",
    size: "sm",
    tone: "success",
  },
  argTypes: {
    tone: {
      control: "select",
      options: ["neutral", "info", "success", "warning", "danger"],
    },
  },
} satisfies Meta<typeof StatusBadge>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const Tones: Story = {
  render: () => (
    <HStack gap="2" flexWrap="wrap">
      {(["neutral", "info", "success", "warning", "danger"] as const).map(
        (tone: StatusBadgeTone) => (
          <StatusBadge key={tone} tone={tone} size="sm">
            {tone}
          </StatusBadge>
        ),
      )}
    </HStack>
  ),
};

export const StatusVersusLabel: Story = {
  render: () => (
    <LStack gap="3">
      <HStack gap="2">
        <StatusBadge tone="success" size="sm">
          Published
        </StatusBadge>
        <StatusBadge tone="warning" size="sm">
          Needs attention
        </StatusBadge>
      </HStack>
      <HStack gap="2">
        <Badge>Web development</Badge>
        <Badge>Announcement</Badge>
      </HStack>
    </LStack>
  ),
  parameters: {
    docs: {
      description: {
        story:
          "StatusBadge communicates lifecycle or execution state. Badge labels categories and metadata. Similar compact shapes do not make these roles interchangeable.",
      },
    },
  },
};
