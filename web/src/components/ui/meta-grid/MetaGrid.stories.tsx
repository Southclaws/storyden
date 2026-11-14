import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { MetaGrid, MetaItem } from ".";

const meta = {
  title: "Components/Data Display/Meta Grid",
  component: MetaGrid,
  parameters: {
    docs: {
      description: {
        component:
          "A description list for bounded object metadata. Use it for labelled dates, identifiers and status details; use PropertyList for compact aligned label/value rows and Table for comparing multiple objects. Default columns fit the available container width and long values wrap. Ordinary page sections remain on the Canvas.",
      },
    },
  },
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

export const NarrowContainers: Story = {
  render: () => (
    <LStack gap="6">
      {(["48", "80", "full"] as const).map((width) => (
        <LStack key={width} width={width} maxW="full">
          <MetaGrid>
            <MetaItem label="Saved">November 1, 2026</MetaItem>
            <MetaItem label="Last checked">Not checked yet</MetaItem>
            <MetaItem label="Metadata retrieved">Not recorded</MetaItem>
            <MetaItem label="Long value">
              a-very-long-generated-identifier-that-should-wrap-without-breaking-layout
            </MetaItem>
          </MetaGrid>
        </LStack>
      ))}
    </LStack>
  ),
};
