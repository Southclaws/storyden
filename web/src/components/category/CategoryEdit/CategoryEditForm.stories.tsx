import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { LStack } from "@/styled-system/jsx";

import { CategoryEditForm } from "./CategoryEditForm";

const source = `data:image/svg+xml,${encodeURIComponent(
  '<svg xmlns="http://www.w3.org/2000/svg" width="1600" height="900"><rect width="1600" height="900" fill="#155e75"/><circle cx="800" cy="450" r="280" fill="#fbbf24"/><rect x="700" y="350" width="200" height="200" fill="#164e63"/></svg>',
)}`;

const meta = {
  title: "Screens/Category Editor",
  component: CategoryEditForm,
  decorators: [
    (Story) => (
      <LStack maxW="2xl">
        <Story />
      </LStack>
    ),
  ],
  parameters: {
    docs: {
      description: {
        component:
          "Category cover editing uses the shared image cropper with a fixed 4:1 frame. Replacement, repositioning and removal remain local until Save; Cancel discards them. Unlike Library main images, category covers have no layout or crop-ratio choices because the same image is used in social previews.",
      },
    },
  },
  args: {
    category: {
      id: "general",
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
      name: "General",
      slug: "general",
      description: "Community conversation",
      colour: "#8577ce",
      sort: 0,
      children: [],
      postCount: 0,
    },
  },
} satisfies Meta<typeof CategoryEditForm>;

export default meta;

type Story = StoryObj<typeof meta>;

export const WithoutCover: Story = {};

export const WithCover: Story = {
  args: {
    category: {
      ...meta.args.category,
      cover_image: {
        id: "cover",
        path: source,
        filename: "cover.svg",
        mime_type: "image/svg+xml",
        width: 1600,
        height: 900,
      },
    },
  },
};
