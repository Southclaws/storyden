import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { ReactList } from ".";

const meta = {
  title: "Components/Data Display/React List",
  component: ReactList,
  parameters: {
    docs: {
      description: {
        component:
          "Use ReactList for emoji reaction counts and an authenticated add-reaction picker. Your own reactions have a selected appearance; guest reactions remain visible but disabled. Pending reactions preserve selection, and count changes respect reduced motion. Use IconButton for a standalone action without reaction counts. Theme hooks: reaction__root, reaction__trigger, reaction__count, and reaction__picker.",
      },
    },
  },
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
  parameters: {
    docs: {
      description: {
        story:
          "Guest reactions retain the same fill as unselected authenticated reactions while remaining disabled.",
      },
    },
  },
  args: {
    isLoggedIn: false,
    reactions: meta.args.reactions.map((reaction) => ({
      ...reaction,
      hasReacted: false,
    })),
  },
};

export const Pending: Story = {
  args: {
    pendingEmojis: new Set(["🔥"]),
  },
};

export const Empty: Story = {
  args: {
    reactions: [],
  },
};
