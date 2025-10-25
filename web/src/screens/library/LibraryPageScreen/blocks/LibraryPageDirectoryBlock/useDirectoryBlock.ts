import { keyBy } from "lodash";
import { useEffect } from "react";

import { LibraryPageBlockTypeDirectory } from "@/lib/library/metadata";
import { deepEqual } from "@/utils/equality";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useBlock } from "../useBlock";

import { MappableNodeFields, getDefaultBlockConfig } from "./column";

export function useDirectoryBlock(): LibraryPageBlockTypeDirectory {
  const { store } = useLibraryPageContext();
  const block = useBlock("directory");
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  if (block === undefined) {
    throw new Error(
      "useDirectoryBlock rendered in a page without a Directory block.",
    );
  }

  const { overwriteBlock } = store.getState();

  const defaultConfig = getDefaultBlockConfig(currentChildPropertySchema);

  // Self-heal directory block config:
  // If the block config is empty, set it to the default state using schema.
  // If the columns contain any fields that do not exist on the schema, fix it.
  useEffect(() => {
    const timeout = setTimeout(() => {
      if (block?.config === undefined) {
        const defaultBlockConfig = getDefaultBlockConfig(
          currentChildPropertySchema,
        );

        console.debug(
          "Directory block config is undefined, setting to default",
          defaultBlockConfig,
        );

        overwriteBlock({
          type: "directory",
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

        const fixedFields = MappableNodeFields.filter(
          (fixed) => columnFieldsByKey[`fixed:${fixed}`] === undefined,
        ).map((fixed) => ({
          fid: `fixed:${fixed}`,
          hidden: true,
        }));

        const updatedColumnConfig = [
          ...fixedFields,
          ...filteredRemovedColumns,
          ...filteredAddedColumns,
        ];

        if (deepEqual(updatedColumnConfig, block.config.columns)) {
          return;
        }

        console.debug(
          "Directory block config mismatches schema, setting to new config",
          updatedColumnConfig,
        );

        overwriteBlock({
          type: "directory",
          config: {
            layout: block.config.layout,
            columns: updatedColumnConfig,
          },
        });
      }
    }, 0);

    return () => {
      clearTimeout(timeout);
    };
  }, [block]);

  if (block.config === undefined) {
    return {
      ...block,
      config: defaultConfig,
    };
  }

  return block;
}
