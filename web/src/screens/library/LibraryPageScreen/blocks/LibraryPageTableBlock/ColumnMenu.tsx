import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { EyeIcon, FilterIcon } from "lucide-react";
import { PropsWithChildren } from "react";

import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import { LibraryPageBlockTypeTable } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";

import { ColumnDefinition, getDefaultBlockConfig } from "./column";

type Props = {
  column: ColumnDefinition;
};

export function ColumnMenu({ column, children }: PropsWithChildren<Props>) {
  const { node, form } = useLibraryPageContext();

  const currentMetadata = form.watch("meta", node.meta);
  const currentChildPropertySchema = form.watch(
    "childPropertySchema",
    node.child_property_schema,
  );

  const currentTableBlockIndex = currentMetadata.layout?.blocks.findIndex(
    (b) => b.type === "table",
  );
  if (!currentTableBlockIndex) {
    console.warn(
      "attempting to render a ColumnMenu without a table block in the form metadata",
    );
    return null;
  }
  const currentTableBlock = currentMetadata.layout?.blocks[
    currentTableBlockIndex
  ] as LibraryPageBlockTypeTable;

  if (currentTableBlock.config === undefined) {
    currentTableBlock.config = getDefaultBlockConfig(
      currentChildPropertySchema,
    );
  }

  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "hide-show": {
        handleColumnHide();
        break;
      }
    }
  }

  function handleColumnNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    const nextChildPropertySchema = currentChildPropertySchema.map((ps) => {
      if (ps.fid === column.fid) {
        return {
          ...ps,
          name: event.target.value,
        };
      }

      return ps;
    });

    form.setValue("childPropertySchema", nextChildPropertySchema);
  }

  function handleColumnHide() {
    const nextBlocks =
      currentMetadata.layout?.blocks.map((block) => {
        if (block.type === "table") {
          const nextColumns =
            currentTableBlock.config?.columns.map((col) => {
              if (col.fid === column.fid) {
                return { ...col, hidden: true };
              }
              return col;
            }) ?? [];

          const nextTable = {
            ...block,
            config: {
              ...block.config,
              columns: nextColumns,
            } as LibraryPageBlockTypeTable["config"],
          };

          return nextTable;
        }

        return block;
      }) ?? [];

    form.setValue("meta.layout.blocks", nextBlocks);
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect} size="xs">
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            {!column.fixed && (
              <Menu.ItemGroup pl="2" py="2">
                <Input
                  size="sm"
                  value={column.name}
                  onChange={handleColumnNameChange}
                />
              </Menu.ItemGroup>
            )}

            <Menu.ItemGroup>
              <Menu.Item value="hide-show">
                <EyeIcon />
                &nbsp;Hide
              </Menu.Item>

              {!column.fixed && (
                <Menu.Item value="delete">
                  <DeleteIcon />
                  &nbsp;Delete
                </Menu.Item>
              )}

              <Menu.Item value="filter">
                <FilterIcon />
                &nbsp;Filter...
              </Menu.Item>
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
