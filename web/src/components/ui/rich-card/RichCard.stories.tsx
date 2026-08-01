import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Badge } from "../badge";
import { MoreIcon } from "../icons/More";

import { Card, CardGrid, CardRows } from ".";

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
  args: {
    id: "example-card",
    title: "Example card",
    url: "/d/example-card",
  },
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

export const Variants: Story = {
  render: () => (
    <div style={{ display: "grid", gap: "1rem", maxWidth: "48rem" }}>
      {(["row", "responsive", "box", "fill"] as const).map((shape) => (
        <Card
          key={shape}
          id={`card-${shape}`}
          title={`Card shape ${shape}`}
          url={`/d/card-${shape}`}
          text="A shared card primitive used for discussions, links, and library previews."
          image="https://images.unsplash.com/photo-1518005020951-eccb494ad742?auto=format&fit=crop&w=900&q=80"
          shape={shape}
          disableAnchors
        >
          <Badge>{shape}</Badge>
        </Card>
      ))}
      {(["default", "emphasized", "accent"] as const).map((backgroundColor) => (
        <Card
          key={backgroundColor}
          id={`card-background-${backgroundColor}`}
          title={`Background ${backgroundColor}`}
          url={`/d/card-background-${backgroundColor}`}
          text="Background variants keep the same structure while changing surface emphasis."
          backgroundColor={backgroundColor}
          menu={<MoreIcon />}
          disableAnchors
        >
          <Badge>{backgroundColor}</Badge>
        </Card>
      ))}
    </div>
  ),
};
