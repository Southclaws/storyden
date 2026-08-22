import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

import { RelativeTime } from ".";

const meta = {
  title: "Components/Data Display/Relative Time",
  component: RelativeTime,
  parameters: {
    docs: {
      description: {
        component:
          "Renders an accessible relative timestamp with approximate or strict wording. Use it instead of formatting relative dates directly at product call sites. Pair it with a separate absolute timestamp when the exact occurrence matters; it does not replace domain-specific date formatting.",
      },
    },
  },
  args: {
    value: new Date(Date.now() - 90 * 60 * 1000),
    precision: "approximate",
  },
  argTypes: {
    precision: {
      control: "select",
      options: ["approximate", "strict"],
    },
  },
} satisfies Meta<typeof RelativeTime>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {};

export const WithExactOccurrence: Story = {
  render: () => {
    const value = new Date(Date.now() - 90 * 60 * 1000);

    return (
      <LStack gap="1">
        <RelativeTime value={value} fontWeight="semibold" />
        <Text variant="metadata">Scheduled run at 22 Aug 2026, 18:00</Text>
      </LStack>
    );
  },
  parameters: {
    docs: {
      description: {
        story:
          "Use relative time as the scan-friendly label and retain the precise domain timestamp as supporting context when users may need to compare or audit occurrences.",
      },
    },
  },
};
