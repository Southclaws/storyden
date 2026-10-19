import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { ComponentProps } from "react";

import { LStack, styled } from "@/styled-system/jsx";
import { richCard } from "@/styled-system/recipes";
import { pluralise } from "@/utils/text";

import { Badge } from "../badge";
import { IconButton } from "../icon-button";
import { MoreIcon } from "../icons/More";
import { PinIcon } from "../icons/Pin";
import { ToolsetIcon } from "../icons/Toolset";
import { SectionHeading } from "../section-heading";
import { Text } from "../text";

import { Card, CardGrid, CardRows } from ".";

const shapes = richCard.variantMap.shape;
const backgrounds = richCard.variantMap.backgroundColor;
const previewImage = new URL(
  "../../../stories/assets/mountains.png",
  import.meta.url,
).href;
const previewArgs = {
  id: "alpine-notebook",
  title: "Alpine field notes",
  url: "/l/alpine-field-notes",
  text: "A quiet lake, warm evening light, and notes from our first walk through the valley.",
  image: previewImage,
  imageAlt: "Sunlit mountain peaks reflected in a still lake",
  header: <Text variant="metadata">Environment art</Text>,
  children: <Text variant="metadata">Mara · 12 replies</Text>,
  controls: <Badge>Field guide</Badge>,
  menu: (
    <IconButton aria-label="More options" type="button" variant="ghost">
      <MoreIcon />
    </IconButton>
  ),
  disableAnchors: true,
} satisfies ComponentProps<typeof Card>;

const richContent = [
  "<h3>Along the shoreline</h3>",
  ...Array.from(
    { length: 8 },
    (_, index) =>
      `<p>Field note ${index + 1}: We followed the water around the valley, recording the changing light on the mountain slopes and the reflections along the shore.</p>`,
  ),
].join("");

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
  argTypes: {
    shape: { control: "select", options: shapes },
    backgroundColor: { control: "select", options: backgrounds },
    headingAs: { control: "select", options: ["h2", "h3", "h4"] },
    image: { control: "text" },
    imageAlt: { control: "text" },
    titleStatus: { control: "text" },
    content: { control: "text" },
    header: { control: false },
    menu: { control: false },
    children: { control: false },
    controls: { control: false },
  },
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component:
          "The canonical Surface family for independent content objects. Every Surface and Card has a `border.default` boundary; fill alone is not enough separation from Canvas. Use `Card` when content has a title, description, metadata, media, or controls; use `CardRows` and `CardGrid` for repeated collections. `titleIcon` is an optional, decorative 16px semantic marker whose size and muted colour are owned by Card. Use a product-owned icon when it adds deliberate object identity or state; do not substitute a visually similar icon from another domain. Card exposes a named article and an h2 title by default; use headingAs for a nested section. The title is the primary link and the excerpt remains readable text. Use titleStatus for meaningful state conveyed by a decorative icon. Backdrop images are decorative; foreground previews default to decorative teasers, with imageAlt available when a meaningful image description is supplied. Do not copy the title into imageAlt. The recipe has four shapes and three backgrounds, with no size axis. Responsive layout changes at a container width of 560px. The fill shape intentionally hides its title slot and top-right menu: put visible identity in its header; its content provides the primary link. Use Variants and Without Media to inspect all 12 combinations, and Container Widths to check 320, 480, 560, 561, and 800px frames. Do not wrap arbitrary page sections in cards, and use an Inset treatment for subordinate content inside a Surface.",
      },
    },
  },
} satisfies Meta<typeof Card>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = { args: previewArgs };

export const Row: Story = {
  args: {
    ...previewArgs,
    id: "example-row",
    title: "Elastic build machines are now available",
    url: "/d/example-row",
    text: "A discussion preview with enough text to exercise the row layout, title rhythm, and footer content without depending on API data.",
    shape: "row",
    children: <Badge>Announcement</Badge>,
    disableAnchors: true,
  },
};

export const ObjectIcon: Story = {
  args: {
    id: "plugin-builder-toolset",
    title: "Plugin Studio",
    url: "/robots/toolsets/plugin-builder",
    text: "Builds, tests, packages, and validates Storyden plugins.",
    titleIcon: <ToolsetIcon />,
    shape: "row",
    disableAnchors: true,
  },
  parameters: {
    docs: {
      description: {
        story:
          "A recurring product object may use its dedicated semantic identity icon. `ToolsetIcon` is backed by Lucide Toolbox; it does not borrow `CollectionIcon` merely because both metaphors contain grouped objects.",
      },
    },
  },
};

export const StateIcon: Story = {
  args: {
    id: "pinned-discussion",
    title: "Pinned community guidelines",
    url: "/d/pinned-discussion",
    text: "A title icon may communicate exceptional state. Card owns its compact size, alignment, and colour. The status text also identifies the state to screen readers.",
    titleIcon: <PinIcon />,
    titleStatus: "Pinned discussion",
    shape: "row",
    disableAnchors: true,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Use a semantic icon export when it adds deliberate product identity or state. This example communicates pinned state; object-family icons should likewise have their own product name rather than borrowing a visually similar symbol.",
      },
    },
  },
};

export const Rows: Story = {
  args: {
    id: "example-rows",
    url: "/d/example-rows",
  },
  render: () => <CardRows items={exampleItems} />,
};

export const ContentPrefetch: Story = {
  args: {
    id: "content-prefetch",
    title: "Community handbook",
    url: "/l/community-handbook",
    text: "Hover or focus the title link to warm this page once.",
    prefetchOnIntent: true,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Thread and Library previews opt into prefetchOnIntent. Nothing is prefetched on mount: hover or keyboard focus warms the destination once per card. Other cards, disabled anchors, and external URLs do not prefetch. Use this for content readers are likely to open, rather than profile, settings, or administrative destinations.",
      },
    },
  },
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
              {count} {pluralise(count, "card")}
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

export const Responsive: Story = {
  args: { ...previewArgs, id: "responsive-card", shape: "responsive" },
};

export const Box: Story = {
  args: { ...previewArgs, id: "box-card", shape: "box" },
};

export const Fill: Story = {
  args: {
    ...previewArgs,
    id: "fill-card",
    shape: "fill",
    header: <Text fontWeight="semibold">Alpine field notes</Text>,
  },
  parameters: {
    docs: {
      description: {
        story:
          "The image fills the card. The header carries visible identity, while content and footer sit over the lower backdrop. The title and top-right menu slots are intentionally hidden by this recipe. The content link remains keyboard-accessible.",
      },
    },
  },
};

export const Variants: Story = {
  args: previewArgs,
  argTypes: { shape: { control: false }, backgroundColor: { control: false } },
  render: (args) => (
    <LStack gap="8" maxWidth="3xl">
      {shapes.map((shape) => (
        <LStack key={shape} gap="3" as="section">
          <SectionHeading>{shape}</SectionHeading>
          {backgrounds.map((backgroundColor) => (
            <LStack key={backgroundColor} gap="1">
              <Text variant="metadata">{backgroundColor}</Text>
              <Card
                {...args}
                id={`variant-${shape}-${backgroundColor}`}
                shape={shape}
                backgroundColor={backgroundColor}
                headingAs="h3"
                header={
                  shape === "fill" ? (
                    <Text fontWeight="semibold">Alpine field notes</Text>
                  ) : (
                    args.header
                  )
                }
              />
            </LStack>
          ))}
        </LStack>
      ))}
    </LStack>
  ),
};

export const WithoutMedia: Story = {
  ...Variants,
  args: { ...previewArgs, image: undefined },
  parameters: {
    docs: {
      description: {
        story:
          "All shape/background combinations without media. This exposes reserved tracks and the empty fill treatment instead of concealing them with a photograph.",
      },
    },
  },
};

export const ContainerWidths: Story = {
  args: previewArgs,
  argTypes: { shape: { control: false } },
  render: (args) => (
    <LStack gap="8" width="full">
      {shapes.map((shape) => (
        <LStack key={shape} as="section" gap="3" width="full">
          <SectionHeading>{shape}</SectionHeading>
          {[320, 480, 560, 561, 800].map((width) => (
            <LStack key={width} gap="1" width="full">
              <Text variant="metadata">
                {width}px container
                {width === 560
                  ? " · last narrow width"
                  : width === 561
                    ? " · first wide width"
                    : ""}
              </Text>
              <styled.div width="full" overflowX="auto">
                <div style={{ width }}>
                  <Card
                    {...args}
                    id={`width-${shape}-${width}`}
                    shape={shape}
                    headingAs="h3"
                    header={
                      shape === "fill" ? (
                        <Text fontWeight="semibold">Alpine field notes</Text>
                      ) : (
                        args.header
                      )
                    }
                  />
                </div>
              </styled.div>
            </LStack>
          ))}
        </LStack>
      ))}
    </LStack>
  ),
  parameters: {
    docs: {
      description: {
        story:
          "Every shape in exact-width frames, including both sides of the 560px container query. Scroll within a frame when the canvas is narrower; the card's test width stays exact. Use the background control to inspect each surface treatment at these widths.",
      },
    },
  },
};

export const HeaderAndMenu: Story = {
  args: { ...previewArgs, shape: "responsive" },
  render: (args) => (
    <LStack gap="6" maxWidth="3xl">
      {[
        { name: "Header and menu", header: true, menu: true },
        { name: "Header only", header: true, menu: false },
        { name: "Menu only", header: false, menu: true },
        { name: "Neither header nor menu", header: false, menu: false },
      ].map((state, index) => (
        <LStack key={state.name} as="section" gap="2">
          <SectionHeading>{state.name}</SectionHeading>
          <Card
            {...args}
            id={`slots-${index}`}
            headingAs="h3"
            header={state.header ? args.header : undefined}
            menu={state.menu ? args.menu : undefined}
          />
        </LStack>
      ))}
    </LStack>
  ),
  parameters: {
    docs: {
      description: {
        story:
          "Use the shape and background controls to check slot placement. Row and responsive menus change placement when the header is absent. In narrow responsive cards, the menu uses the content column so media keeps its full inset width, with or without a header. Fill deliberately hides its top-right menu.",
      },
    },
  },
};

export const ContentStates: Story = {
  args: { ...previewArgs, shape: "row" },
  render: (args) => (
    <LStack gap="6" maxWidth="3xl">
      {[
        {
          name: "Title only",
          props: {
            text: undefined,
            content: undefined,
            image: undefined,
            header: undefined,
            menu: undefined,
            children: undefined,
            controls: undefined,
          },
        },
        {
          name: "Long title and unbroken text",
          props: {
            title:
              "An unusually long field notebook title about light, mountain reflections, and finding a safe route around the valley",
            text: "https://example.com/field-notebooks/a-very-long-unbroken-reference-that-must-not-widen-the-card-or-push-the-actions-outside-the-frame",
          },
        },
        {
          name: "Missing title",
          props: {
            title: undefined,
            text: "An untitled entry still has readable content and a way to open it.",
          },
        },
        { name: "Missing description", props: { text: undefined } },
        {
          name: "Long rich content",
          props: { text: undefined, content: richContent },
        },
      ].map((state, index) => (
        <LStack key={state.name} as="section" gap="2">
          <SectionHeading>{state.name}</SectionHeading>
          <Card
            {...args}
            {...state.props}
            id={`content-${index}`}
            headingAs="h3"
          />
        </LStack>
      ))}
    </LStack>
  ),
  parameters: {
    docs: {
      description: {
        story:
          "Check sparse content, long strings, untitled objects, and rich-content disclosure in every shape using the controls. Tab to Show more and use Enter or Space; the button should retain focus while the preview expands and collapses.",
      },
    },
  },
};

export const AccessiblePreview: Story = {
  args: {
    id: "accessible-preview",
    title: "A walk around the harbour",
    url: "/t/harbour",
    text: "Notes on the new walkway and the boats arriving this spring.",
    image: previewImage,
    imageAlt: "Sunlit mountain peaks reflected in a still lake",
    shape: "responsive",
  },
  parameters: {
    docs: {
      description: {
        story:
          "Use heading navigation to scan titles, then read the excerpt without an extra link. A supplied foreground description is announced once; the blurred duplicate is hidden. Tab reaches the title and any independent controls. Long rich content exposes a Show more button that retains focus and reports its expanded state.",
      },
    },
  },
};
