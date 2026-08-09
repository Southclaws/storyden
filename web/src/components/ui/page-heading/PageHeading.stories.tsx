import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { PageHeading } from ".";

const meta = {
  title: "Components/Typography/Page Heading",
  component: PageHeading,
  parameters: {
    docs: {
      description: {
        component:
          "The screen identity. PageHeading always renders `h1`, appears at most once per screen, and is normally composed through PageHeader. Surface titles and authored-content headings use their own systems rather than changing this component's size or element.",
      },
    },
  },
  args: {
    children: "System settings",
  },
} satisfies Meta<typeof PageHeading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};
