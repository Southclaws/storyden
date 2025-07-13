"use client";

import { SelectValueChangeDetails, createListCollection } from "@ark-ui/react";
import { useState } from "react";

import { Node } from "@/api/openapi-schema";
import { LibraryPageSelect } from "@/components/library/LibraryPageSelect";
import { CancelAction } from "@/components/site/Action/Cancel";
import { EditAction } from "@/components/site/Action/Edit";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { useFeedContext } from "@/screens/feed/FeedContext";
import { HStack, styled } from "@/styled-system/jsx";

const sources = [
  {
    label: "Threads",
    value: "threads" as const,
    icon: <DiscussionIcon width="4" />,
  },
  {
    label: "Library",
    value: "library" as const,
    icon: <LibraryIcon width="4" />,
  },
  {
    label: "Categories",
    value: "categories" as const,
    icon: <CategoryIcon width="4" />,
  },
];

export function FeedConfig() {
  const { isEditingEnabled, isEditing, feed, updateFeed, handleToggleEditing } =
    useFeedContext();

  if (!isEditingEnabled) {
    return null;
  }

  const collection = createListCollection({ items: sources });

  async function handleHomepageNodeChange(node: Node | undefined) {
    await updateFeed({
      layout: {
        type: "list",
      },
      source: {
        type: "library",
        node: node?.id,
      },
    });
  }

  async function handleSourceTypeChange({ value }: SelectValueChangeDetails) {
    if (value.length === 0) {
      return;
    }

    const feedSourceType = value[0] as typeof feed.source.type;

    await updateFeed({
      layout: {
        type: "list",
      },
      source: {
        type: feedSourceType,
      },
    });
  }

  return (
    <HStack w="full" justify="end">
      {isEditing ? (
        <>
          {feed.source.type === "library" && (
            <>
              <LibraryPageSelect
                onChange={handleHomepageNodeChange}
                value={feed.source.node}
              />
            </>
          )}

          <Select.Root
            w="fit"
            size="xs"
            collection={collection}
            defaultValue={[feed.source.type]}
            positioning={{ sameWidth: false }}
            onValueChange={handleSourceTypeChange}
          >
            <Select.Control>
              <Select.Trigger>
                <Select.ValueText placeholder="Select a Source" />
                <SelectIcon />
              </Select.Trigger>
            </Select.Control>
            <Select.Positioner>
              <Select.Content>
                {sources.map((item) => (
                  <Select.Item key={item.value} item={item}>
                    <Select.ItemText mr="2">
                      <HStack gap="1">
                        <styled.span w="4">{item.icon}</styled.span>
                        <styled.span>{item.label}</styled.span>
                      </HStack>
                    </Select.ItemText>
                    <Select.ItemIndicator>
                      <CheckIcon />
                    </Select.ItemIndicator>
                  </Select.Item>
                ))}
              </Select.Content>
            </Select.Positioner>
          </Select.Root>

          <CancelAction onClick={handleToggleEditing}>Done</CancelAction>
        </>
      ) : (
        <EditAction onClick={handleToggleEditing}>Configure feed</EditAction>
      )}
    </HStack>
  );
}
