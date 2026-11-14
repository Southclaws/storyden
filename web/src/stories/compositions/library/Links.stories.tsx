import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useRouter } from "@storybook/nextjs-vite/navigation.mock";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { fn } from "storybook/test";

import type { Asset, Link, ProfileReference } from "@/api/openapi-schema";
import {
  LinkCardGrid,
  LinkCardRows,
} from "@/components/library/links/LinkCardList";
import { LinkDetails } from "@/components/library/links/LinkDetails";
import { LStack } from "@/styled-system/jsx";

const image: Asset = {
  id: "field-notes-image",
  filename: "mountains.png",
  path: new URL("../../assets/mountains.png", import.meta.url).href,
  mime_type: "image/png",
  width: 1200,
  height: 700,
};

const author: ProfileReference = {
  id: "curator",
  handle: "curator",
  name: "Curator",
  joined: "2026-01-01T00:00:00Z",
  roles: [],
};

const saved: Link = {
  id: "field-notes",
  slug: "field-notes",
  url: "https://fieldnotes.example/landscapes",
  domain: "fieldnotes.example",
  title: "Landscape field notes",
  description:
    "Light, colour and composition studies from a week in the mountains.",
  createdAt: "2026-11-01T09:00:00Z",
  updatedAt: "2026-11-13T14:00:00Z",
  liveness_status: "alive",
  liveness_checked_at: "2026-11-13T14:00:00Z",
  scrape_status: "success",
  scrape_attempted_at: "2026-11-12T10:00:00Z",
  scraped_at: "2026-11-12T10:00:00Z",
  primary_image: image,
  assets: [image],
  nodes: [
    {
      id: "environment-art",
      slug: "environment-art",
      name: "Environment art",
      description: "References, studies and notes shared by the community.",
      createdAt: "2026-10-01T09:00:00Z",
      updatedAt: "2026-11-13T14:00:00Z",
      owner: author,
      assets: [],
      hide_child_tree: false,
      tags: [],
      visibility: "published",
      meta: {},
    },
  ],
  posts: [
    {
      id: "landscape-study",
      slug: "landscape-study",
      title: "What are you studying this week?",
      description: "Share a reference and what you learned from it.",
      createdAt: "2026-11-03T09:00:00Z",
      updatedAt: "2026-11-03T09:00:00Z",
      author,
      visibility: "published",
      assets: [],
      reacts: [],
      collections: { in_collections: 0, has_collected: false },
      likes: { likes: 0, liked: false },
    },
  ],
  recomentations: [],
};

const pending: Link = {
  ...saved,
  id: "pending",
  slug: "pending",
  domain: "new-reference.example",
  url: "https://new-reference.example/",
  title: undefined,
  description: undefined,
  primary_image: undefined,
  assets: [],
  nodes: [],
  posts: [],
  liveness_status: "unknown",
  liveness_checked_at: undefined,
  scrape_status: "pending",
  scrape_attempted_at: undefined,
  scraped_at: undefined,
};

const olderPreview: Link = {
  ...saved,
  id: "older-preview",
  slug: "older-preview",
  title: "An older saved reference",
  liveness_status: "unknown",
  liveness_checked_at: undefined,
  scrape_status: "pending",
  scrape_attempted_at: undefined,
  scraped_at: undefined,
};

const unavailable: Link = {
  ...saved,
  id: "unavailable",
  slug: "unavailable",
  title: "A saved reference from an unreachable site",
  liveness_status: "offline",
  scrape_status: "failed",
};

const failedRefresh: Link = {
  ...saved,
  id: "failed-refresh",
  slug: "failed-refresh",
  title: "A preview kept after a failed refresh",
  scrape_status: "failed",
};

const withoutImage: Link = {
  ...saved,
  id: "without-image",
  slug: "without-image",
  title: "Reading list: colour and contrast",
  primary_image: undefined,
  assets: [],
};

const longTitle: Link = {
  ...withoutImage,
  id: "long-title",
  slug: "long-title",
  title:
    "A practical collection of landscape references, colour studies, composition notes and lighting experiments for environment artists",
};

const longDomain = [
  "a".repeat(63),
  "b".repeat(63),
  "c".repeat(63),
  "d".repeat(53),
  "example",
].join(".");

const longDomainLink: Link = {
  ...withoutImage,
  id: "long-domain",
  slug: "long-domain",
  title: "Long domain reference",
  domain: longDomain,
  url: `https://${longDomain}/reference`,
  scrape_status: "failed",
};

const longDomainWithoutTitle: Link = {
  ...longDomainLink,
  id: "long-domain-without-title",
  slug: "long-domain-without-title",
  title: undefined,
};

const links = [
  saved,
  pending,
  unavailable,
  failedRefresh,
  withoutImage,
  longTitle,
  olderPreview,
  longDomainLink,
  longDomainWithoutTitle,
];

const meta = {
  title: "Compositions/Library/Links",
  parameters: {
    nextjs: { appDirectory: false },
    docs: {
      description: {
        component:
          "Browse-first Links experiments using the canonical Surface, StatusBadge and metadata components. All content and health results are illustrative fixtures. Rows and Grid use identical records. OlderSavedPreview represents existing installations with saved metadata but no recorded scrape or check timestamps; pending must not hide that preview or imply active fetching.",
      },
    },
  },
  decorators: [
    (Story) => (
      <NuqsTestingAdapter>
        <LStack w="full" maxW="6xl" mx="auto">
          <Story />
        </LStack>
      </NuqsTestingAdapter>
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
      bfcacheId: "link-stories",
    });
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const Rows: Story = {
  render: () => <LinkCardRows links={links} />,
};

export const Grid: Story = {
  render: () => <LinkCardGrid links={links} />,
};

export const NarrowGrid: Story = {
  render: () => (
    <LStack w="80" maxW="full">
      <LinkCardGrid links={[longDomainLink, longDomainWithoutTitle, pending]} />
    </LStack>
  ),
};

export const Detail: Story = {
  render: () => <LinkDetails link={saved} />,
};

export const DetailImageOnly: Story = {
  render: () => <LinkDetails link={{ ...saved, description: undefined }} />,
};

export const Unreachable: Story = {
  render: () => <LinkDetails link={unavailable} />,
};

export const Pending: Story = {
  render: () => <LinkDetails link={pending} />,
};

export const OlderSavedPreview: Story = {
  render: () => <LinkDetails link={olderPreview} />,
};

export const LongTitle: Story = {
  render: () => <LinkDetails link={longTitle} />,
};
