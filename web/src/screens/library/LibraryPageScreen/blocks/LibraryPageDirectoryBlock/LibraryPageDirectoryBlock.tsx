import { ChangeEvent } from "react";

import { NodeTree } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { Input } from "@/components/ui/input";
import { Box, HStack, LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { AddPropertyMenu } from "./AddPropertyMenu/AddPropertyMenu";
import {
  LibraryPageDirectoryBlockContextProvider,
  useDirectoryBlockContext,
} from "./Context";
import { LibraryPageDirectoryBlockGrid } from "./LibraryPageDirectoryBlockGrid";
import { LibraryPageDirectoryBlockTable } from "./LibraryPageDirectoryBlockTable";
import { PropertyListMenu } from "./PropertyListMenu/PropertyListMenu";
import { useChildrenWithTags } from "./useChildrenWithTags";
import { useDirectoryBlock } from "./useDirectoryBlock";

export function LibraryPageDirectoryBlock() {
  return (
    <LibraryPageDirectoryBlockContextProvider>
      <LibraryPageDirectoryBlockContents />
    </LibraryPageDirectoryBlockContextProvider>
  );
}

export function LibraryPageDirectoryBlockContents() {
  const {
    handleSearch,
    handleTagFilter,
    highlightedTags,
    searchQuery,
    tagFilters,
    childrenSort,
  } = useDirectoryBlockContext();
  const { nodeID, initialChildren } = useLibraryPageContext();
  const { editing } = useEditState();

  function handleSearchChange(event: ChangeEvent<HTMLInputElement>) {
    handleSearch(event.target.value);
  }

  const { data, hasChildren, error, tags } = useChildrenWithTags(
    nodeID,
    initialChildren,
    childrenSort,
    tagFilters,
    searchQuery,
  );

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LStack w="full" gap="2">
      <WStack bgColor="bg.subtle" borderRadius="sm" p="1">
        <TagBadgeList
          tags={tags}
          type="button"
          onClick={handleTagFilter}
          highlightedTags={highlightedTags}
        />
        <HStack gap="1" alignItems="end" flexWrap="wrap-reverse" justify="end">
          <Input
            variant="ghost"
            placeholder="Search..."
            size="xs"
            onChange={handleSearchChange}
            minW="20"
            maxW="min"
            flexShrink="1"
          />

          {editing && (
            <HStack gap="1">
              <AddPropertyMenu unavailable={!hasChildren}>
                <IconButton
                  size="xs"
                  variant="ghost"
                  title="Add a new property."
                >
                  <AddIcon />
                </IconButton>
              </AddPropertyMenu>

              <PropertyListMenu>
                <IconButton size="xs" variant="ghost">
                  <MenuIcon />
                </IconButton>
              </PropertyListMenu>
            </HStack>
          )}
        </HStack>
      </WStack>

      <LibraryPageDirectoryBlockLayout nodes={data.nodes} />
    </LStack>
  );
}

function LibraryPageDirectoryBlockLayout({ nodes }: { nodes: NodeTree }) {
  const block = useDirectoryBlock();

  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  if (!block) {
    console.warn(
      "attempting to render a LibraryPageDirectoryBlock without a block in the page metadata",
    );
    return null;
  }

  // Get layout from config, default to table
  const layout = block.config?.layout ?? "table";

  // Switch between different layout views
  switch (layout) {
    case "grid":
      return (
        <LibraryPageDirectoryBlockGrid
          nodes={nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
        />
      );

    case "table":
      return (
        <LibraryPageDirectoryBlockTable
          nodes={nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
        />
      );

    default:
      return <Unready error={new Error(`Unknown layout type: "${layout}"`)} />;
  }
}
