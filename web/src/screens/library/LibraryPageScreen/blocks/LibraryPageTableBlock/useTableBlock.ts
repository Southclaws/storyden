import { LibraryPageBlockTypeTable } from "@/lib/library/metadata";

import { useWatch } from "../../store";
import { useBlock } from "../useBlock";

import { getDefaultBlockConfig } from "./column";

export function useTableBlock(): LibraryPageBlockTypeTable {
  const block = useBlock("table");
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  const defaultConfig = getDefaultBlockConfig(currentChildPropertySchema);

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
