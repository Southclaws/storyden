import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";

import { Breadcrumbs } from ".";

const meta = {
  title: "Components/Navigation/Breadcrumbs",
  component: Breadcrumbs,
  parameters: {
    docs: {
      description: {
        component:
          "Describes the current page's place in a Library, discussion, or other information hierarchy. Breadcrumbs answer 'where am I?'; use BackAction when the requirement is 'return to where I came from'. Keep labels compact and allow deep trailing items to truncate.",
      },
    },
  },
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
