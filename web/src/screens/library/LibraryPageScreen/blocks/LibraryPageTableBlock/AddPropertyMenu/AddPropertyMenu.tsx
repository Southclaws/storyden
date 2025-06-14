import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { useClickAway } from "@uidotdev/usehooks";
import { PropsWithChildren, useState } from "react";

import { handle } from "@/api/client";
import { nodeUpdateChildrenPropertySchema } from "@/api/openapi-client/nodes";
import { PropertyType } from "@/api/openapi-schema";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { useTableBlock } from "../useTableBlock";

export function AddPropertyMenu({ children }: PropsWithChildren) {
  const { currentNode, store } = useLibraryPageContext();
  const [name, setName] = useState<string>("");
  const { addChildProperty } = store.getState();

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
        currentNode.id,
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

    addChildProperty(newProperty);

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
