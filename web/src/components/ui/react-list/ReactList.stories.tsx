import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { ReactList } from ".";

const meta = {
  title: "Components/Data Display/React List",
  component: ReactList,
  args: {
    isLoggedIn: true,
    reactions: [
      {
        emoji: "🔥",
        count: 3,
        hasReacted: true,
        reactedBy: ["odin", "southclaws", "molly"],
      },
      {
        emoji: "👍",
        count: 1,
        hasReacted: false,
        reactedBy: ["odin"],
      },
    ],
    onReactExisting: () => undefined,
    onReactPicker: () => undefined,
  },
} satisfies Meta<typeof ReactList>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Guest: Story = {
  args: {
    isLoggedIn: false,
  },
};

export const Pending: Story = {
  args: {
    pendingEmojis: new Set(["🔥"]),
  },
};
