import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import { Breadcrumbs } from ".";

const meta = {
  title: "UI/Breadcrumbs",
  component: Breadcrumbs,
} satisfies Meta<typeof Breadcrumbs>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    index: { label: "Library", href: "/l" },
    crumbs: [
      { label: "Documentation Hub", href: "/l/docs" },
      { label: "Advanced Topics", href: "/l/docs/advanced" },
      {
        label: "Very long page title that should truncate gracefully",
        href: "/l/docs/advanced/long-title",
      },
    ],
    children: (
      <Button size="sm" variant="subtle">
        New page
      </Button>
    ),
  },
};
