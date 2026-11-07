import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useId, useState } from "react";

import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { LStack } from "@/styled-system/jsx";

import {
  ToggleMatrix,
  type ToggleMatrixColumn,
  type ToggleMatrixRow,
  type ToggleMatrixValue,
  type ToggleMatrixValueChangeDetails,
} from ".";

const meta = {
  title: "Components/Forms/Toggle Matrix",
  component: ToggleMatrix,
  parameters: {
    docs: {
      description: {
        component:
          "Compares boolean settings across two axes. Each cell is a toggle button with its row and column in its accessible name. Tab moves between cells; Space or Enter toggles them. Use Toggle Group for a single axis and Radio Group for mutually exclusive choices. Relationships between settings belong to the consumer.",
      },
    },
  },
} satisfies Meta<typeof ToggleMatrix>;

export default meta;

type Story = StoryObj<typeof NotificationExample>;

const rows: ToggleMatrixRow[] = [
  {
    value: "replies",
    label: "Replies",
    description: "Replies to your threads and replies to you.",
  },
  {
    value: "library",
    label: "Library contributions",
    description: "Suggested changes to pages you own.",
  },
  {
    value: "trails",
    label: "Trail alerts",
    description: "Trail runs that need your attention.",
  },
  {
    value: "reports",
    label: "Reports and moderation",
    description: "Reports and moderation updates relevant to you.",
  },
];

const columns: ToggleMatrixColumn[] = [
  { value: "inbox", label: "Inbox" },
  { value: "push", label: "Push" },
];

const initialValue: ToggleMatrixValue = {
  replies: { inbox: true, push: true },
  library: { inbox: true, push: false },
  trails: { inbox: true, push: true },
  reports: { inbox: true, push: false },
};

function NotificationExample({
  pushDisabled = false,
  disabled = false,
  mixed = false,
}: {
  pushDisabled?: boolean;
  disabled?: boolean;
  mixed?: boolean;
}) {
  const descriptionId = useId();
  const [value, setValue] = useState<ToggleMatrixValue>({
    ...initialValue,
    library: { inbox: true, push: mixed ? "mixed" : false },
  });

  function handleValueChange({
    value: nextValue,
    row,
    column,
    pressed,
  }: ToggleMatrixValueChangeDetails) {
    const selected = { ...nextValue[row] };

    if (column === "push" && pressed) {
      selected["inbox"] = true;
    }

    if (column === "inbox" && !pressed) {
      selected["push"] = false;
    }

    setValue({ ...nextValue, [row]: selected });
  }

  const visibleColumns = columns.map((column) => ({
    ...column,
    disabled: column.value === "push" && pushDisabled,
  }));

  return (
    <LStack gap="3" width="full" maxWidth="xl">
      <LStack gap="1">
        <SectionHeading>Notification delivery</SectionHeading>
        <Text variant="supporting" id={descriptionId}>
          Choose where you receive each notification. Push also enables Inbox;
          turning Inbox off turns Push off.
        </Text>
      </LStack>
      <ToggleMatrix
        label="Notification delivery"
        aria-describedby={descriptionId}
        rowHeading="Notify me about"
        rows={rows}
        columns={visibleColumns}
        value={value}
        onValueChange={handleValueChange}
        disabled={disabled}
      />
    </LStack>
  );
}

export const Notifications: Story = {
  render: () => <NotificationExample />,
};

export const PushUnavailable: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Disabling a delivery column preserves its selections while the remaining controls can still be changed.",
      },
    },
  },
  render: () => <NotificationExample pushDisabled />,
};

export const Disabled: Story = {
  render: () => <NotificationExample disabled />,
};

export const MixedEvents: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "A filled cell with a dash means some underlying settings are enabled. Activating it enables all settings represented by that cell.",
      },
    },
  },
  render: () => <NotificationExample mixed />,
};
