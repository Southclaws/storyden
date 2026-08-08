import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { MetaGrid, MetaItem } from ".";

const meta = {
  title: "Components/Data Display/Meta Grid",
  component: MetaGrid,
  args: {
    children: null,
  },
} satisfies Meta<typeof MetaGrid>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <MetaGrid>
      <MetaItem label="Created">August 1, 2026</MetaItem>
      <MetaItem label="Author">@southclaws</MetaItem>
      <MetaItem label="Status">Published</MetaItem>
      <MetaItem label="Long value">
        a-very-long-generated-identifier-that-should-wrap-without-breaking-layout
      </MetaItem>
    </MetaGrid>
  ),
};
