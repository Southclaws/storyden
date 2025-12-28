import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren, useState } from "react";

import { handle } from "@/api/client";
import { nodeUpdateChildrenPropertySchema } from "@/api/openapi-client/nodes";
import { PropertyType } from "@/api/openapi-schema";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import { styled } from "@/styled-system/jsx";
import { useClickAway } from "@/utils/useClickAway";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { useDirectoryBlock } from "../useDirectoryBlock";

type AddPropertyMenuProps = PropsWithChildren<{
  unavailable?: boolean;
}>;

export function AddPropertyMenu({
  children,
  unavailable = false,
}: AddPropertyMenuProps) {
  const { nodeID, store } = useLibraryPageContext();
  const [name, setName] = useState<string>("");
  const { addChildProperty } = store.getState();
  const directoryBlock = useDirectoryBlock();

  // Menu opening logic. We circumvent the default behaviour here because we
  // want to keep the menu open during the creation of a new property then close
  // it as soon as the change has been committed to the node's state store.
  const [open, setOpen] = useState(false);
  const ref = useClickAway<HTMLDivElement>(() => setOpen(false));

  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  function handleColumnNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    setName(event.target.value);
  }

  async function handleSave(): Promise<boolean> {
    const trimmed = name.trim();
    if (trimmed === "") {
      return false;
    }

    const exists = currentChildPropertySchema.find((p) => p.name === trimmed);
    if (exists) {
      const column = directoryBlock.config?.columns.find(
        (c) => c.fid === exists.fid,
      );
      if (column?.hidden) {
        throw new Error(
          `Property "${trimmed}" already exists but is hidden in the directory. Open the directory menu to toggle the column's visibility.`,
        );
      } else {
        throw new Error(`Property "${trimmed}" already exists.`);
      }
    }

    const updatedChildPropertySchema = [
      ...currentChildPropertySchema,
      {
        name: trimmed,
        sort: "",
        type: PropertyType.text,
      },
    ];

    const newSchema = await nodeUpdateChildrenPropertySchema(
      nodeID,
      updatedChildPropertySchema,
    );
    if (!newSchema) {
      throw new Error("Failed to update page properties");
    }

    const newProperty = newSchema.properties.find((p) => p.name === trimmed);
    if (!newProperty) {
      throw new Error("New property not found in updated schema");
    }

    addChildProperty(newProperty);

    setName("");

    return true;
  }

  async function handleSelect(value: MenuSelectionDetails) {
    if (value.value === "create") {
      await handle(async () => {
        const close = await handleSave();
        if (close) {
          setOpen(() => false);
        }
      });
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
              <Input
                size="sm"
                value={name}
                onChange={handleColumnNameChange}
                disabled={unavailable}
              />
            </Menu.ItemGroup>

            <Menu.ItemGroup pl="2" py="1">
              <Menu.Item value="create" disabled={unavailable}>
                Create
              </Menu.Item>

              {unavailable && (
                <Menu.ItemGroup>
                  <styled.p fontSize="xs" color="fg.muted">
                    Properties require at least one sub-page.
                  </styled.p>
                </Menu.ItemGroup>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
