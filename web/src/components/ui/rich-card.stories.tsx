import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Badge } from "./badge";
import { Card, CardGrid, CardRows } from "./rich-card";

const exampleItems = [
  {
    id: "storyden-cli",
    title: "Storyden command line interface",
    url: "/d/storyden-cli",
    text: "Ship, inspect, and operate a Storyden community from the terminal with a workflow designed for local-first maintainers.",
  },
  {
    id: "library-updates",
    title: "Library updates for community docs",
    url: "/d/library-updates",
    text: "A calmer reading surface for documentation hubs, guides, drafts, and long-lived knowledge shared by the community.",
  },
  {
    id: "moderation-flows",
    title: "Moderation flows",
    url: "/d/moderation-flows",
    text: "Practical queue and review patterns for keeping a community healthy without turning every screen into an admin console.",
  },
];

const meta = {
  title: "UI/Rich Card",
  component: Card,
  parameters: {
    layout: "padded",
  },
} satisfies Meta<typeof Card>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Row: Story = {
  args: {
    id: "example-row",
    title: "Elastic build machines are now available",
    url: "/d/example-row",
    text: "A discussion preview with enough text to exercise the row layout, title rhythm, and footer content without depending on API data.",
    shape: "row",
    children: <Badge>Announcement</Badge>,
    disableAnchors: true,
  },
};

export const Rows: Story = {
  args: {
    id: "example-rows",
    url: "/d/example-rows",
  },
  render: () => <CardRows items={exampleItems} />,
};

export const Grid: Story = {
  args: {
    id: "example-grid",
    url: "/d/example-grid",
  },
  render: () => <CardGrid items={exampleItems} />,
};
