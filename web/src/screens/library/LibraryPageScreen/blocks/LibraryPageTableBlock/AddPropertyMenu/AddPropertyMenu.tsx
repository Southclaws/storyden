import { Portal } from "@ark-ui/react";
import { PropsWithChildren, useState } from "react";

import { PropertySchema, PropertySchemaList } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeTable,
  LibraryPageBlockTypeTableColumn,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../../Context";
import { getDefaultBlockConfig } from "../column";

export function AddPropertyMenu({ children }: PropsWithChildren) {
  const { node, form } = useLibraryPageContext();
  const [name, setName] = useState<string>("");

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

  const currentBlocks = [...(currentMetadata.layout?.blocks || [])];

  function handleColumnNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    setName(event.target.value);
  }

  function handleSave() {
    const trimmed = name.trim();
    if (trimmed === "") {
      return;
    }

    const newColumn: LibraryPageBlockTypeTableColumn = {
      fid: `new_field_${Date.now()}`,
      hidden: false,
    };

    const nextBlocks = currentBlocks.map((block) => {
      if (block.type === "table") {
        return {
          ...block,
          config: {
            ...block.config,
            columns: [...(block.config?.columns ?? []), newColumn],
          },
        };
      }
      return block;
    });

    const updatedMeta = {
      ...currentMetadata,
      layout: {
        ...currentMetadata.layout,
        blocks: nextBlocks,
      },
    };

    console.log(updatedMeta);

    const updatedChildPropertySchema: PropertySchemaList = [
      ...currentChildPropertySchema,
      {
        fid: newColumn.fid,
        name: trimmed,
        sort: "",
        type: "text",
      } satisfies PropertySchema,
    ];

    form.setValue("meta", updatedMeta);
    form.setValue("childPropertySchema", updatedChildPropertySchema);

    setName(""); // Reset the input field
  }

  return (
    <Menu.Root
      closeOnSelect={true}
      lazyMount
      size="xs"
      positioning={{
        placement: "bottom-end",
      }}
    >
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup pl="2" py="1">
              <Menu.ItemGroupLabel>New property</Menu.ItemGroupLabel>
              <Input size="sm" value={name} onChange={handleColumnNameChange} />
            </Menu.ItemGroup>

            {/* <Menu.ItemGroup>
              <Menu.Item value="hide-show">
                <EyeIcon />
                &nbsp;Hide
              </Menu.Item>
            </Menu.ItemGroup> */}

            <Menu.ItemGroup pl="2" py="1">
              <Button
                size="xs"
                variant="subtle"
                type="button"
                onClick={handleSave}
              >
                Create
              </Button>
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
