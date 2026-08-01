import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { DatePicker, DateRangePicker } from ".";

const meta = {
  title: "UI/Date Picker",
  component: DatePicker,
} satisfies Meta<typeof DatePicker>;

export default meta;

type Story = StoryObj<typeof meta>;

export const SingleDate: Story = {
  render: () => <DatePicker />,
};

export const DateRange: Story = {
  render: () => (
    <LStack gap="4" maxW="md">
      <DateRangePicker />
      <DateRangePicker triggerLabel="Filter by date range" hideInputs />
    </LStack>
  ),
};
