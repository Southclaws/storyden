import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box, LStack } from "@/styled-system/jsx";

import { DatePicker, DateRangePicker } from ".";

const meta = {
  title: "Components/Forms/Date Picker",
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

export const AtViewportEdge: Story = {
  render: () => (
    <Box display="flex" justifyContent="end" minH="md" width="full">
      <DateRangePicker triggerLabel="Filter by date range" hideInputs />
    </Box>
  ),
};
