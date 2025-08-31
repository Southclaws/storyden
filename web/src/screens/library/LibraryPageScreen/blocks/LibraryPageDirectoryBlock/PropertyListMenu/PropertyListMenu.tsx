import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { HideIcon } from "@/components/ui/icons/HideIcon";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import * as Menu from "@/components/ui/menu";
import { BlockIcon } from "@/lib/library/blockIcons";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { mergeFieldsAndPropertySchema } from "../column";
import { useDirectoryBlock } from "../useDirectoryBlock";

type Props = {
  hideFixedFields?: boolean;
};

export function PropertyListMenu({
  children,
  hideFixedFields = false,
}: PropsWithChildren<Props>) {
  const { store } = useLibraryPageContext();
  const { setChildPropertyHiddenState } = store.getState();

  const currentDirectoryBlock = useDirectoryBlock();
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  function handleSelect(value: MenuSelectionDetails) {
    const fid = value.value;

    const hidden =
      currentDirectoryBlock.config?.columns?.find((c) => c.fid === fid)
        ?.hidden ?? true;

    setChildPropertyHiddenState(fid, !hidden);
  }

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    currentDirectoryBlock,
    true,
  )
    .filter((c) => (hideFixedFields ? !c._fixedFieldName : true))
    // NOTE: Primary image is handled separately. It's only present in the grid
    // view mode. Plus, while it's *technically* a "field", to a user it's not
    // really a field, it's a separate feature of a page unrelated to properties.
    .filter((c) => c.fid !== "fixed:primary_image");

  const supportsCoverImage = currentDirectoryBlock.config?.layout === "grid";
  const coverImageHiddenState =
    currentDirectoryBlock.config?.columns.find(
      (c) => c.fid === "fixed:primary_image",
    )?.hidden ?? false;

  return (
    <Menu.Root
      lazyMount
      size="xs"
      positioning={{
        placement: "bottom-end",
      }}
      closeOnSelect={false}
      onSelect={handleSelect}
    >
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36" maxW="max">
            {supportsCoverImage && (
              <Menu.ItemGroup pl="2" py="1">
                <Menu.ItemGroupLabel>Options</Menu.ItemGroupLabel>
                <Menu.Item value="fixed:primary_image">
                  <HStack w="full">
                    <HStack w="full" gap="1" textWrap="nowrap">
                      <BlockIcon blockType="cover" />
                      <span>Cover image</span>
                    </HStack>

                    {coverImageHiddenState ? <HideIcon /> : <ShowIcon />}
                  </HStack>
                </Menu.Item>
              </Menu.ItemGroup>
            )}

            <Menu.ItemGroup pl="2" py="1">
              <Menu.ItemGroupLabel>Properties</Menu.ItemGroupLabel>

              {columns.map((property) => (
                <Menu.Item key={property.fid} value={property.fid}>
                  {property.hidden ? <HideIcon /> : <ShowIcon />}
                  &nbsp;
                  {property.name}
                </Menu.Item>
              ))}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
