import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { DefaultSettings } from "@/lib/settings/settings";
import { styled } from "@/styled-system/jsx";

import { BrandSettingsForm } from "./BrandSettings";

const settings = {
  ...DefaultSettings,
  motd: {
    content: "<p>Community maintenance starts tonight at 21:00.</p>",
    start_at: "2030-01-15T21:00:00.000Z",
    end_at: "2030-01-16T00:00:00.000Z",
    metadata: { type: "information" as const },
  },
};

const meta = {
  title: "Screens/Admin/Brand Settings",
  component: BrandSettingsForm,
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta<typeof BrandSettingsForm>;

export default meta;

type Story = StoryObj<typeof meta>;

export const WithMessageOfTheDay: Story = {
  args: { settings },
  render: (args) => (
    <styled.main marginX="auto" maxW="5xl" padding={{ base: "4", md: "8" }}>
      <BrandSettingsForm {...args} />
    </styled.main>
  ),
};
