import { dequal } from "dequal";
import { keyBy } from "lodash";
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

  const { overwriteBlock } = store.getState();

  const defaultConfig = getDefaultBlockConfig(currentChildPropertySchema);

  // Self-heal table block config:
  // If the block config is empty, set it to the default state using schema.
  // If the columns contain any fields that do not exist on the schema, fix it.
  useEffect(() => {
    const timeout = setTimeout(() => {
      if (block?.config === undefined) {
        const defaultBlockConfig = getDefaultBlockConfig(
          currentChildPropertySchema,
        );

        console.debug(
          "Table block config is undefined, setting to default",
          defaultBlockConfig,
        );

        overwriteBlock({
          type: "table",
          config: defaultBlockConfig,
        });
      } else {
        const schemaFieldsByKey = keyBy(currentChildPropertySchema, "fid");
        const columnFieldsByKey = keyBy(block.config.columns, "fid");

        // First, remove columns that do not exist in the schema any more.
        const filteredRemovedColumns = block.config.columns.filter(
          (c) => c.fid in schemaFieldsByKey || c.fid.startsWith("fixed:"),
        );

        // Then, add any columns from the schema that aren't in the column config.
        const filteredAddedColumns = currentChildPropertySchema
          .filter((p) => !(p.fid in columnFieldsByKey))
          .map((p) => ({
            fid: p.fid,
            hidden: false,
          }));

        const updatedColumnConfig = [
          ...filteredRemovedColumns,
          ...filteredAddedColumns,
        ];

        if (dequal(updatedColumnConfig, block.config.columns)) {
          return;
        }

        console.debug(
          "Table block config mismatches schema, setting to new config",
          updatedColumnConfig,
        );

        overwriteBlock({
          type: "table",
          config: {
            columns: updatedColumnConfig,
          },
        });
      }
    }, 0);

    return () => {
      clearTimeout(timeout);
    };
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
