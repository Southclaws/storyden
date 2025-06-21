import { useEffect } from "react";

import { LibraryPageBlockTypeTable } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useBlock } from "../useBlock";

import { getDefaultBlockConfig } from "./column";

export function useTableBlock(): LibraryPageBlockTypeTable {
  const { store } = useLibraryPageContext();
  const block = useBlock("table");
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  if (block === undefined) {
    throw new Error("useTableBlock rendered in a page without a Table block.");
  }

  const defaultConfig = getDefaultBlockConfig(currentChildPropertySchema);

  // Self-heal table block config:
  // If the block config is empty, set it to the default state using schema.
  // If the columns contain any fields that do not exist on the schema, fix it.
  useEffect(() => {
    if (block?.config === undefined) {
      const defaultBlockConfig = getDefaultBlockConfig(
        currentChildPropertySchema,
      );

      store.getState().overwriteBlock({
        type: "table",
        config: defaultBlockConfig,
      });
    } else {
    }
  }, [block]);

  if (block === undefined) {
    return {
      type: "table",
      config: defaultConfig,
    };
  }

  if (block.config === undefined) {
    return {
      ...block,
      config: defaultConfig,
    };
  }

  return block;
}
