import { useNodeListChildren } from "@/api/openapi-client/nodes";
import { useSortIndicator } from "@/components/site/SortIndicator";
import { Unready } from "@/components/site/Unready";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";

import { LibraryPageDirectoryBlockGrid } from "./LibraryPageDirectoryBlockGrid";
import { LibraryPageDirectoryBlockTable } from "./LibraryPageDirectoryBlockTable";
import { useDirectoryBlock } from "./useDirectoryBlock";

export function LibraryPageDirectoryBlock() {
  const { nodeID, initialChildren, store } = useLibraryPageContext();
  const { sort, handleSort } = useSortIndicator();

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

  const nodes = data.nodes;

  if (nodes.length === 0) {
    return null;
  }

  if (!block) {
    console.warn(
      "attempting to render a LibraryPageDirectoryBlock without a block in the form metadata",
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
    default:
      return (
        <LibraryPageDirectoryBlockTable
          nodes={nodes}
          block={block}
          currentChildPropertySchema={currentChildPropertySchema}
          sort={sort}
          handleSort={handleSort}
        />
      );
  }
}
