import { useNodeListChildren } from "@/api/openapi-client/nodes";
import { EmptyState } from "@/components/site/EmptyState";
import { useSortIndicator } from "@/components/site/SortIndicator";
import { Unready } from "@/components/site/Unready";
import { Center } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

import {
  LibraryPageDirectoryBlockContextProvider,
  useDirectoryBlockContext,
} from "./Context";
import { LibraryPageDirectoryBlockGrid } from "./LibraryPageDirectoryBlockGrid";
import { LibraryPageDirectoryBlockTable } from "./LibraryPageDirectoryBlockTable";
import { useDirectoryBlock } from "./useDirectoryBlock";

export function LibraryPageDirectoryBlock() {
  return (
    <LibraryPageDirectoryBlockContextProvider>
      <LibraryPageDirectoryBlockContents />
    </LibraryPageDirectoryBlockContextProvider>
  );
}

export function LibraryPageDirectoryBlockContents() {
  const { nodeID, initialChildren } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();
  const { searchQuery } = useDirectoryBlockContext();

  // format the sort property as "name" or "-name" for asc/desc
  const childrenSort =
    sort !== null
      ? sort?.order === "asc"
        ? sort.property
        : `-${sort.property}`
      : undefined;

  const { data, error } = useNodeListChildren(
    nodeID,
    {
      children_sort: childrenSort,
      q: searchQuery || undefined,
    },
    {
      swr: {
        fallbackData: initialChildren,
      },
    },
  );

  const block = useDirectoryBlock();
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  if (!data) {
    return <Unready error={error} />;
  }

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
          nodes={data.nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
        />
      );

    case "table":
      return (
        <LibraryPageDirectoryBlockTable
          nodes={data.nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
          sort={sort}
          handleSort={handleSort}
        />
      );

    default:
      return <Unready error={new Error(`Unknown layout type: "${layout}"`)} />;
  }
}
