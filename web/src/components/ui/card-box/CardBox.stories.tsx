import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import { CardBox } from ".";

const meta = {
  title: "Components/Data Display/Card Box",
  component: CardBox,
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "A compact compatibility Surface for independent arbitrary content. It always retains the canonical Surface border, padding, and radius. Prefer the richer Surface components for titled, repeated, media-backed, or metadata-heavy objects; do not use CardBox as a generic wrapper around top-level settings or page sections.",
      },
    },
  },
} satisfies Meta<typeof CardBox>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Variants: Story = {
  render: () => (
    <styled.div display="grid" gap="4" width="96">
      <CardBox>
        <styled.strong fontSize="sm">Discussion summary</styled.strong>
        <styled.span color="text.muted" fontSize="xs">
          A compact independent object on the application canvas.
        </styled.span>
      </CardBox>

      <CardBox kind="edge">
        <styled.div padding="2">
          <styled.strong fontSize="sm">Edge content</styled.strong>
        </styled.div>
        <styled.div
          backgroundColor="background.inset"
          borderTopColor="border.default"
          borderTopWidth="thin"
          color="text.muted"
          fontSize="xs"
          padding="2"
        >
          The edge variant delegates internal spacing to its content.
        </styled.div>
      </CardBox>
    </styled.div>
  ),
};
