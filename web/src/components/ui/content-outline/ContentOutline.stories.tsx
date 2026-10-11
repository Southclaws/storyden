import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box, styled } from "@/styled-system/jsx";

import {
  ContentOutline,
  type ContentOutlineEntry,
  collectContentHeadings,
} from ".";

const meta = {
  title: "Components/Navigation/ContentOutline",
  component: ContentOutline,
  args: { label: "On this page", collectEntries: collect },
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "A desktop reading aid for headings, page blocks and current-page replies. The quiet tick rail becomes prominent on hover or focus and opens a navigable outline. Overflowing ends fade as a scroll cue, and a small separator distinguishes replies from the original content. Reply rows pair author and relative date above a snippet. Screens provide rendered targets through a scoped collector. Use SectionNavigation for route navigation, TreeView for hierarchy browsing and Popover for generic supporting content. Do not use while editing, for embedded previews or on mobile. Keyboard users can press Enter to enter the outline; Escape restores focus and selecting a row focuses its destination.",
      },
    },
  },
} satisfies Meta<typeof ContentOutline>;
export default meta;
type Story = StoryObj<typeof meta>;

const copy =
  "Good environment art gives the viewer somewhere to look and a reason to stay. Collect references, block out the main shapes, and use lighting to make the important parts easy to read. Leave room for small details that suggest how people use the place.";

function collect(root: HTMLElement): ContentOutlineEntry[] {
  const sections = root.querySelectorAll<HTMLElement>("section");
  const entries: ContentOutlineEntry[] = [];

  for (const [index, target] of sections.entries()) {
    const key = String(index);
    const isReply = target.hasAttribute("data-reply");

    if (isReply) {
      const createdAt = new Date(Date.now() - index * 60 * 60 * 1000);

      entries.push({
        key,
        kind: "reply",
        label: `Maker ${index}`,
        description: copy,
        createdAt: createdAt.toISOString(),
        depth: 0,
        target,
      });

      continue;
    }

    const label = target.querySelector("h2")?.textContent;

    if (!label) {
      continue;
    }

    entries.push({ key, kind: "section", label, depth: 0, target });

    const headings = collectContentHeadings(target, key).slice(1);

    entries.push(...headings);
  }

  return entries;
}

function ReadingExample({
  count = 5,
  replies = false,
  disabled = false,
}: {
  count?: number;
  replies?: boolean;
  disabled?: boolean;
}) {
  return (
    <Box maxW="5xl" mx="auto" p="6">
      <ContentOutline
        label={replies ? "In this thread" : "On this page"}
        collectEntries={collect}
        disabled={disabled}
      >
        <styled.h1 fontSize="2xl" mb="6">
          Environment art field notes
        </styled.h1>
        {Array.from({ length: count }, (_, index) => (
          <styled.section
            key={index}
            data-reply={replies && index > 0 ? "" : undefined}
            pb="12"
          >
            <styled.h2 fontSize="xl" mb="3">
              {index === 0
                ? "Finding the story"
                : `Section ${index + 1}: Building a place that feels inhabited`}
            </styled.h2>
            <p>{copy}</p>
            <styled.h3 fontSize="lg" mt="6" mb="3">
              Details and observations
            </styled.h3>
            <p>{copy}</p>
          </styled.section>
        ))}
      </ContentOutline>
    </Box>
  );
}

export const Page: Story = { render: () => <ReadingExample /> };
export const Thread: Story = { render: () => <ReadingExample replies /> };
export const LongOutline: Story = {
  render: () => <ReadingExample count={80} />,
};
export const LongThread: Story = {
  render: () => <ReadingExample count={80} replies />,
};
export const Editing: Story = { render: () => <ReadingExample disabled /> };
export const Empty: Story = { render: () => <ReadingExample count={0} /> };
