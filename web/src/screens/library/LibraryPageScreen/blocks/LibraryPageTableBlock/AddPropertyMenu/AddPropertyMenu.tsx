import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { useClickAway } from "@uidotdev/usehooks";
import { PropsWithChildren, useState } from "react";

import { handle } from "@/api/client";
import { nodeUpdateChildrenPropertySchema } from "@/api/openapi-client/nodes";
import { PropertyType } from "@/api/openapi-schema";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeTable,
  LibraryPageBlockTypeTableColumn,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { getDefaultBlockConfig } from "../column";

export function AddPropertyMenu({ children }: PropsWithChildren) {
  const { currentNode, store } = useLibraryPageContext();
  const [name, setName] = useState<string>("");
  const { setMeta } = store.getState();

  // Menu opening logic. We circumvent the default behaviour here because we
  // want to keep the menu open during the creation of a new property then close
  // it as soon as the change has been committed to the node's state store.
  const [open, setOpen] = useState(false);
  const ref = useClickAway<HTMLDivElement>(() => setOpen(false));

  const currentMetadata = useWatch((s) => s.draft.meta);
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
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

  async function handleSave(): Promise<boolean> {
    const trimmed = name.trim();
    if (trimmed === "") {
      return false;
    }

    const updatedChildPropertySchema = [
      ...currentChildPropertySchema,
      {
        name: trimmed,
        sort: "",
        type: PropertyType.text,
      },
    ];

    const newSchema = await handle(async () => {
      return await nodeUpdateChildrenPropertySchema(
        // NOTE: Will break if slug changes - replace with live watch state.
        currentNode.slug,
        updatedChildPropertySchema,
      );
    });
    if (!newSchema) {
      console.error("Failed to update child property schema");
      return false;
    }

    const newProperty = newSchema.properties.find((p) => p.name === trimmed);
    if (!newProperty) {
      console.error("New property not found in updated schema");
      return false;
    }

    const newColumn: LibraryPageBlockTypeTableColumn = {
      fid: newProperty.fid,
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

    setMeta(updatedMeta);

    setName("");

    return true;
  }

  async function handleSelect(value: MenuSelectionDetails) {
    if (value.value === "create") {
      const close = await handleSave();
      if (close) {
        setOpen(() => false);
      }
    }
  }

  return (
    <Menu.Root
      open={open}
      lazyMount
      size="xs"
      positioning={{
        placement: "bottom-end",
      }}
      closeOnSelect={false}
      onSelect={handleSelect}
      onEscapeKeyDown={() => setOpen(false)}
    >
      <Menu.Trigger asChild onClick={() => setOpen(!open)}>
        {children}
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner ref={ref}>
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
              <Menu.Item
                value="create"
                bgColor="bg.subtle"
                _hover={{
                  bgColor: "bg.muted",
                }}
              >
                Create
              </Menu.Item>
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
