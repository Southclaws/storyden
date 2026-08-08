import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { CSSProperties } from "react";

import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import * as Table from "@/components/ui/table";
import { Text } from "@/components/ui/text";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

const directoryRows = [
  ["API reference", "Reference", "Updated today"],
  ["Installation", "Guide", "Updated yesterday"],
  ["Network configuration", "Reference", "Updated 4 days ago"],
];

const meta = {
  title: "Compositions/Library/Page Blocks",
  parameters: {
    layout: "padded",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

type DirectoryPreviewProps = {
  candidate?: string;
  inputBackground?: CSSProperties["backgroundColor"];
};

function DirectoryPreview({
  candidate,
  inputBackground = "var(--colors-background-inset)",
}: DirectoryPreviewProps) {
  return (
    <LStack data-candidate={candidate} gap="2" width="full">
      <WStack background="background.inset" borderRadius="sm" padding="1">
        <HStack gap="1" flexWrap="wrap">
          <Badge size="sm">documentation</Badge>
          <Badge size="sm">networking</Badge>
        </HStack>
        <Input
          aria-label="Search library pages"
          flexShrink="1"
          maxW="40"
          minW="20"
          px="2"
          placeholder="Search..."
          style={{ backgroundColor: inputBackground }}
          variant="inset"
        />
      </WStack>

      <Box borderRadius="sm" overflow="hidden">
        <Table.Root size="sm" variant="dense">
          <Table.Head>
            <Table.Row>
              <Table.Header>Page</Table.Header>
              <Table.Header>Type</Table.Header>
              <Table.Header>Activity</Table.Header>
            </Table.Row>
          </Table.Head>
          <Table.Body>
            {directoryRows.map(([page, type, activity]) => (
              <Table.Row key={page}>
                <Table.Cell>{page}</Table.Cell>
                <Table.Cell>{type}</Table.Cell>
                <Table.Cell>{activity}</Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      </Box>
    </LStack>
  );
}

export const DirectoryTable: Story = {
  render: () => (
    <styled.main maxW="3xl" width="full">
      <DirectoryPreview />
    </styled.main>
  ),
};

export const DirectorySearchBackgrounds: Story = {
  render: () => (
    <LStack gap="8" maxW="3xl" width="full">
      <LStack gap="2">
        <Text variant="metadata">Control background - neutral step 1</Text>
        <DirectoryPreview
          candidate="control"
          inputBackground="var(--colors-background-control)"
        />
      </LStack>

      <LStack gap="2">
        <Text variant="metadata">Inset control - neutral step 2</Text>
        <DirectoryPreview
          candidate="control-inset"
          inputBackground="var(--colors-background-control-inset)"
        />
      </LStack>

      <LStack gap="2">
        <Text variant="metadata">Inset background - neutral step 3</Text>
        <DirectoryPreview
          candidate="inset"
          inputBackground="var(--colors-background-inset)"
        />
      </LStack>
    </LStack>
  ),
};
