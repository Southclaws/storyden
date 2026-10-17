import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { ReactList } from ".";

const meta = {
  title: "Components/Data Display/React List",
  component: ReactList,
  parameters: {
    docs: {
      description: {
        component:
          "Use ReactList for existing emoji reactions and an authenticated add-reaction picker. Your own reactions use the theme's selected surface and an inset border; aria-pressed exposes the same toggle state to assistive technology. Guest reactions retain their button fill while remaining disabled, and pending reactions preserve selection. Reaction buttons use system emoji characters; the native emoji picker loads when opened. Counts animate with CSS only when the user's motion preference allows it. Use IconButton for a single standalone action rather than reaction counts and selection.",
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
