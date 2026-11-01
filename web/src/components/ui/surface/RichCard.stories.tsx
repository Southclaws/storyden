import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useRouter } from "@storybook/nextjs-vite/navigation.mock";
import type { ComponentProps } from "react";
import { fn } from "storybook/test";

import { useSearchHighlights } from "@/components/search/useSearchHighlights";
import { LStack, styled } from "@/styled-system/jsx";
import { richCard } from "@/styled-system/recipes";
import { pluralise } from "@/utils/text";

import { Badge } from "../badge";
import { IconButton } from "../icon-button";
import { EmptyIcon } from "../icons/Empty";
import { MoreIcon } from "../icons/More";
import { PinIcon } from "../icons/Pin";
import { ToolsetIcon } from "../icons/Toolset";
import { Input } from "../input";
import { PropertyList } from "../property-list";
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
  header: <Badge size="sm">Environment art</Badge>,
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
  // These isolated cards only need navigation spies. The installed App Router
  // story boundary does not support Next's current render-tree context.
  beforeEach: () => {
    useRouter.mockReturnValue({
      back: fn(),
      forward: fn(),
      refresh: fn(),
      push: fn(),
      replace: fn(),
      prefetch: fn(),
      bfcacheId: "surface-story",
    });
  },
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
    nextjs: { appDirectory: false },
    layout: "padded",
    docs: {
      description: {
        component:
          "The canonical Surface family for independent content objects. Row, responsive, box, and fill cards use a `border.default` boundary. Minimal uses a flat inset surface for media-led collections. Use `Card` when content has a title, description, metadata, media, or controls; use `CardRows` and `CardGrid` for repeated collections. `titleIcon` is an optional, decorative 16px semantic marker whose size and muted colour are owned by Card. Use a product-owned icon when it adds deliberate object identity or state; do not substitute a visually similar icon from another domain. Card exposes a named article and an h2 title by default; use headingAs for a nested section. The title is the primary link and the excerpt remains readable text. Use titleStatus for meaningful state conveyed by a decorative icon. Backdrop images are decorative; foreground previews default to decorative teasers, with imageAlt available when a meaningful image description is supplied. Do not copy the title into imageAlt. The recipe has five shapes and three backgrounds, with no size axis. Responsive layout changes at a container width of 560px. The fill shape intentionally hides its title slot and top-right menu: put visible identity in its header; its content provides the primary link. Use Variants and Without Media to inspect all 15 combinations, and Container Widths to check 320, 480, 560, 561, and 800px frames. Do not wrap arbitrary page sections in cards, and use an Inset treatment for subordinate content inside a Surface.",
      },
    },
  },
} satisfies Meta<typeof Card>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = { args: previewArgs };

export const Minimal: Story = {
  args: {
    ...previewArgs,
    shape: "minimal",
    imageAspectRatio: 4 / 3,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Minimal is the flat, media-first card used by Library directories. Use it for compact visual collections; Editorial (box) places its image after the summary. Minimal shares the same title, body, header, footer, and editing slots. This example pairs a compact header badge with a control menu, matching the directory's header treatment. The menu overlays the top-right corner at a fixed inset; it does not vertically center within taller or wrapping headers. Keep header content compact and reserve space on the right when composing it with a menu. imageProps supplies responsive sources, dimensions, and loading hints; imageFallback supplies a placeholder only when the image is absent. The cover links to the page in view mode and becomes non-navigating with disableAnchors.",
      },
    },
  },
};

export const MinimalMediaStates: Story = {
  render: () => (
    <styled.div
      display="grid"
      gridTemplateColumns="repeat(auto-fit, minmax(min(100%, 220px), 1fr))"
      gap="4"
      maxW="5xl"
      alignItems="start"
    >
      <Card
        id="minimal-cover"
        shape="minimal"
        title="Cover and summary"
        url="/l/notes"
        image={previewImage}
        imageAspectRatio={4 / 3}
        text="A fixed 4:3 cover for a Grid collection."
      />
      <Card
        id="minimal-natural"
        shape="minimal"
        title="Natural image"
        url="/l/notes"
        image={previewImage}
        text="Gallery keeps the image's natural proportions."
      />
      <Card
        id="minimal-placeholder"
        shape="minimal"
        title="Missing cover"
        url="/l/notes"
        imageFallback={<EmptyIcon />}
        text="A processing or missing image retains a placeholder."
      />
      <Card
        id="minimal-hidden"
        shape="minimal"
        title="Hidden cover"
        url="/l/notes"
        text="No image or fallback means no media wrapper or reserved space."
      />
      <Card
        id="minimal-image-only"
        shape="minimal"
        label="Image-only notebook"
        url="/l/notes"
        image={previewImage}
        imageAspectRatio={4 / 3}
      />
      <Card
        id="minimal-label-only"
        shape="minimal"
        label="Field notes"
        url="/l/notes"
        imageFallback={<styled.span lineClamp={3}>Field notes</styled.span>}
      />
    </styled.div>
  ),
};

export const MinimalEditing: Story = {
  args: {
    id: "minimal-editing",
    shape: "minimal",
    title: "Alpine field notes",
    url: "/l/notes",
    image: previewImage,
    imageAspectRatio: 4 / 3,
    disableAnchors: true,
    titleContent: (
      <Input
        aria-label="Page name"
        variant="ghost"
        defaultValue="Alpine field notes"
      />
    ),
    body: (
      <LStack gap="1" w="full">
        <Input
          aria-label="Description"
          variant="ghost"
          defaultValue="Notes from the shoreline."
        />
        <PropertyList
          items={[
            {
              id: "status",
              label: "Status",
              value: (
                <Input
                  aria-label="Status"
                  variant="ghost"
                  textAlign="right"
                  defaultValue="Open"
                />
              ),
            },
            { id: "type", label: "Type", value: "Reference" },
          ]}
        />
      </LStack>
    ),
    children: <Badge size="sm">Field guide</Badge>,
  },
};

export const StructuredContent: Story = {
  args: {
    ...previewArgs,
    shape: "box",
    headingAs: "h3",
    imageAspectRatio: 4 / 3,
    titleContent: (
      <Input
        aria-label="Page name"
        defaultValue="Alpine field notes"
        variant="ghost"
        size="sm"
      />
    ),
    body: (
      <Text variant="supporting">
        Structured details stay visible at narrow card widths.
      </Text>
    ),
  },
  parameters: {
    docs: {
      description: {
        story:
          "Use body for structured content alongside the plain-text preview. titleContent replaces the title link for inline editing and fills the available title width; title still supplies the accessible article name. imageAspectRatio gives box media a fixed or asset-derived ratio. label names image-only cards when the title is hidden. These are composition slots, not new visual styles.",
      },
    },
  },
};

export const FooterSpacing: Story = {
  render: () => (
    <LStack gap="6" maxW="3xl">
      {shapes.map((shape) => (
        <LStack key={shape} as="section" gap="2" w="full">
          <SectionHeading>{shape}</SectionHeading>
          <styled.div
            display="grid"
            gridTemplateColumns="repeat(auto-fit, minmax(min(100%, 240px), 1fr))"
            gap="4"
            w="full"
            alignItems="start"
          >
            {[
              { label: "With footer", footer: <Badge size="sm">Serif</Badge> },
              { label: "Without footer", footer: undefined },
            ].map(({ label, footer }) => (
              <Card
                key={label}
                id={`footer-spacing-${shape}-${label}`}
                title={label}
                url="/l/typeface"
                shape={shape}
                image={previewImage}
                body={
                  <PropertyList
                    items={[
                      { id: "license", label: "License", value: "Free" },
                      { id: "usage", label: "Usage", value: "Editorial" },
                    ]}
                  />
                }
                disableAnchors
              >
                {footer}
              </Card>
            ))}
          </styled.div>
        </LStack>
      ))}
    </LStack>
  ),
  parameters: {
    docs: {
      description: {
        story:
          "Surface owns the spacing between sections. A visible footer adds one 2-unit margin above its contents; footer children should not add top padding for this separation. Minimal keeps its bottom content inset only when no footer follows. Compare cards with and without a footer to check for doubled gaps and leftover space.",
      },
    },
  },
};

export const SearchMatches: Story = {
  args: {
    ...previewArgs,
    text: "Field notes from the lake. ".repeat(60),
    controls: <Badge>Field notes metadata stays unhighlighted</Badge>,
  },
  render: function SearchMatches(args) {
    const ref = useSearchHighlights("field lake");
    return (
      <div ref={ref}>
        <Card {...args} />
      </div>
    );
  },
  parameters: {
    docs: {
      description: {
        story:
          "Search highlighting paints literal query words in Card titles and plain-text previews. Metadata, controls, and preview clipping stay unchanged. Select text normally and check both themes. Cards do not highlight themselves: search containers opt in through useSearchHighlights. This is supplementary emphasis, not selected text, keyboard focus, or backend match evidence; unsupported browsers retain the original Card.",
      },
    },
  },
};

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
  parameters: {
    docs: {
      description: {
        story:
          "Media keeps the same inset at the sides and bottom when there is no footer. The optional footer owns the spacing below the image.",
      },
    },
  },
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
          "The image fills the card. A mask on the rounded root keeps Chrome's composited backdrop inside the corners, including at fractional heights and 2× pixel density. The outer container owns the shadow so masking does not cut it off. Media and the lower backdrop share one silhouette. The inset header carries visible identity, while content and footer sit over the lower backdrop. The title and top-right menu slots are intentionally hidden by this recipe. The content link remains keyboard-accessible.",
      },
    },
  },
};

export const FillUnevenContent: Story = {
  args: { ...previewArgs, shape: "fill", children: undefined },
  beforeEach: () => {
    useRouter.mockReturnValue({
      back: fn(),
      forward: fn(),
      refresh: fn(),
      push: fn(),
      replace: fn(),
      prefetch: fn(),
      bfcacheId: "uneven-overlay-story",
    });
  },
  parameters: {
    nextjs: { appDirectory: false },
    docs: {
      description: {
        story:
          "Cards sharing a grid row keep their lower panels attached to the bottom despite wrapped headers or omitted fields. The image track has a 300px minimum and absorbs the row's spare height; empty fields stay omitted.",
      },
    },
  },
  render: (args) => (
    <styled.div
      display="grid"
      gridTemplateColumns="repeat(auto-fit, minmax(min(100%, 200px), 1fr))"
      gap="2"
      maxW="3xl"
    >
      {[
        { title: "Short title", link: true },
        { title: "A longer title that wraps onto a second line", link: true },
        { title: "No external URL", link: false },
      ].map(({ title, link }, index) => (
        <Card
          {...args}
          key={title}
          id={`uneven-overlay-${index}`}
          title={title}
          header={<Text fontWeight="semibold">{title}</Text>}
          text={undefined}
          menu={undefined}
          controls={undefined}
          body={
            <LStack gap="1" w="full" minW="0">
              <Text variant="supporting" lineClamp={2}>
                Field notes and references collected along the shoreline.
              </Text>
              {link && (
                <Text variant="metadata">https://example.com/notes</Text>
              )}
              <PropertyList
                items={[
                  { id: "style", label: "Style", value: "Reference" },
                  { id: "license", label: "License", value: "Open source" },
                  { id: "use", label: "Use", value: "Research" },
                ]}
              />
            </LStack>
          }
        />
      ))}
    </styled.div>
  ),
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
          "All shape/background combinations without media. Row and responsive cards release the image column so text uses the available width. Fill retains its distinct background treatment.",
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
          name: "Empty optional slots",
          props: {
            text: "Empty header, body and footer slots contribute no wrappers or spacing.",
            content: undefined,
            image: undefined,
            header: false,
            body: <></>,
            menu: null,
            children: [],
            controls: false,
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
          "Use heading navigation to scan titles, then read the excerpt without an extra link. A supplied foreground description is announced once; the blurred duplicate is hidden. Tab reaches the title and any independent controls. Long rich content exposes a Show more disclosure styled like an inline text link, sharing document-link colors and hover underline. It retains button semantics, focus, and expanded state, and supports Enter and Space.",
      },
    },
  },
};
