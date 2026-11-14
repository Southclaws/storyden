import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useRouter } from "@storybook/nextjs-vite/navigation.mock";
import { fn } from "storybook/test";

import type { Reply } from "@/api/openapi-schema";
import { RenderedContentProvider } from "@/components/content/ContentRenderer/RenderedContentContext";
import { HighlightedReply } from "@/components/thread/ThreadConclusion";
import { Box } from "@/styled-system/jsx";

const body =
  "<p>We could rotate the meetup between neighbourhoods so more members can join. A short poll would help us choose the first location.</p>";

const reply: Reply = {
  id: "highlighted-reply",
  root_id: "community-meetup",
  root_slug: "community-meetup",
  slug: "community-meetup#highlighted-reply",
  title: "Planning our next community meetup",
  body,
  body_links: [],
  visibility: "published",
  assets: [],
  reacts: [],
  likes: { likes: 3, liked: false },
  collections: { has_collected: false, in_collections: 0 },
  author: {
    id: "alex",
    handle: "alex",
    name: "Alex",
    joined: "2026-01-01T09:00:00Z",
    roles: [],
  },
  createdAt: "2026-11-13T10:00:00Z",
  updatedAt: "2026-11-13T10:00:00Z",
};

const meta = {
  title: "Compositions/Discussion/Highlighted Reply",
  component: HighlightedReply,
  args: { reply, onCurrentPage: true },
  parameters: {
    layout: "padded",
    nextjs: { appDirectory: false },
    docs: {
      description: {
        component:
          "Identifies a reply highlighted within a discussion. When the reply is already on this page, show a compact navigation row rather than an empty reply preview. On another page, retain the content preview and a link to its original position. This composes CardBox, Byline, LinkButton and the existing reply Card; it is not a generic status banner or a replacement for the reply itself.",
      },
    },
  },
  decorators: [
    (Story) => (
      <RenderedContentProvider
        content={{
          [body]: (
            <p>
              We could rotate the meetup between neighbourhoods so more members
              can join. A short poll would help us choose the first location.
            </p>
          ),
        }}
      >
        <Box maxW="5xl" w="full">
          <Story />
        </Box>
      </RenderedContentProvider>
    ),
  ],
  beforeEach: () => {
    useRouter.mockReturnValue({
      back: fn(),
      forward: fn(),
      refresh: fn(),
      push: fn(),
      replace: fn(),
      prefetch: fn(),
      bfcacheId: "highlighted-reply-story",
    });
  },
} satisfies Meta<typeof HighlightedReply>;

export default meta;

type Story = StoryObj<typeof meta>;

export const OnThisPage: Story = {};

export const OnAnotherPage: Story = {
  args: { onCurrentPage: false },
};

export const NarrowWithLongHandle: Story = {
  args: {
    reply: {
      ...reply,
      author: { ...reply.author, handle: "community-meetup-organiser" },
    },
  },
  decorators: [
    (Story) => (
      <Box maxW="80">
        <Story />
      </Box>
    ),
  ],
};
