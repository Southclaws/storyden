import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Badge } from "../badge";
import { IconButton } from "../icon-button";
import { MoreIcon } from "../icons/More";
import { SectionHeading } from "../section-heading";
import { Text } from "../text";

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

const gridLayouts = [
  "1 full-width card",
  "2 cards",
  "3 cards",
  "2 + 2",
  "3 + 2",
  "3 + 3",
  "3 + 2 + 2",
  "3 + 3 + 2",
  "3 + 3 + 3",
];

function createGridItems(count: number) {
  return Array.from({ length: count }, (_, index) => ({
    id: `grid-${count}-${index + 1}`,
    title: `Card ${index + 1}`,
    url: `/d/grid-${count}-${index + 1}`,
    text: "A representative content card for checking balanced grid rows.",
    disableAnchors: true,
  }));
}

const meta = {
  title: "Components/Data Display/Surface",
  component: Card,
  args: {
    id: "example-card",
    title: "Example card",
    url: "/d/example-card",
  },
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component:
          "The canonical Surface family for independent content objects. Every Surface and Card has a `border.default` boundary; fill alone is not enough separation from Canvas. Use `Card` when content has a title, description, metadata, media, or controls; use `CardRows` and `CardGrid` for repeated collections. Do not wrap arbitrary page sections in cards, and use an Inset treatment for subordinate content inside a Surface.",
      },
    },
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

export const GridCounts: Story = {
  args: {
    id: "grid-counts",
    url: "/d/grid-counts",
  },
  render: () => (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "2rem",
        margin: "0 auto",
        maxWidth: "80rem",
      }}
    >
      {Array.from({ length: 9 }, (_, index) => {
        const count = index + 1;

        return (
          <section key={count}>
            <SectionHeading>
              {count} {count === 1 ? "card" : "cards"}
            </SectionHeading>
            <Text variant="supporting" marginBottom="2">
              Balanced as {gridLayouts[index]}.
            </Text>
            <CardGrid items={createGridItems(count)} />
          </section>
        );
      })}
    </div>
  ),
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
          menu={
            <IconButton aria-label="More options" type="button" variant="ghost">
              <MoreIcon />
            </IconButton>
          }
          disableAnchors
        >
          <Badge>{backgroundColor}</Badge>
        </Card>
      ))}
    </div>
  ),
};
