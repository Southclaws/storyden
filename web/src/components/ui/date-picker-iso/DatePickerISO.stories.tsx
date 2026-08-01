import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import { formatISODate, parseISODate } from ".";

const meta = {
  title: "UI/Date Picker ISO",
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const Formatting: Story = {
  render: () => {
    const parsed = parseISODate("2026-08-01");

    return (
      <styled.dl
        display="grid"
        gridTemplateColumns="max-content 1fr"
        gap="2"
        fontSize="sm"
      >
        <styled.dt fontWeight="semibold">Input</styled.dt>
        <styled.dd>2026-08-01</styled.dd>
        <styled.dt fontWeight="semibold">Parsed</styled.dt>
        <styled.dd>{parsed?.toString()}</styled.dd>
        <styled.dt fontWeight="semibold">Formatted</styled.dt>
        <styled.dd>{parsed ? formatISODate(parsed) : "Invalid"}</styled.dd>
      </styled.dl>
    );
  },
};
