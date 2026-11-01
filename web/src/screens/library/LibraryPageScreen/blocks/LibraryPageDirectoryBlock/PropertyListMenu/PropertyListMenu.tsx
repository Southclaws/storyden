import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { HideIcon } from "@/components/ui/icons/HideIcon";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import * as Menu from "@/components/ui/menu";
import { BlockIcon } from "@/lib/library/blockIcons";
import { HStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import {
  CardStyleMenu,
  LayoutMenu,
} from "../LibraryPageDirectoryBlockMenuItems";
import { ColumnDefinition, mergeFieldsAndPropertySchema } from "../column";
import { useDirectoryBlock } from "../useDirectoryBlock";

type Props = {
  hideFixedFields?: boolean;
};

export function PropertyListMenu({
  children,
  hideFixedFields = false,
}: PropsWithChildren<Props>) {
  const { store } = useLibraryPageContext();
  const { setChildPropertyHiddenState, overwriteBlock } = store.getState();

  const currentDirectoryBlock = useDirectoryBlock();
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  function handleSelect(value: MenuSelectionDetails) {
    const fid = value.value;

    if (!currentDirectoryBlock.config?.columns.some((c) => c.fid === fid)) {
      return;
    }

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

  const supportsCoverImage = currentDirectoryBlock.config?.layout !== "table";
  const supportsRichContent =
    currentDirectoryBlock.config?.layout === "rows" ||
    (supportsCoverImage &&
      Boolean(currentDirectoryBlock.config?.card_style) &&
      currentDirectoryBlock.config?.card_style !== "minimal");
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
          <Menu.Content minW="44" maxW="max">
            <LayoutMenu />
            <CardStyleMenu />
            <Menu.Separator />
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
                {supportsRichContent && (
                  <Menu.CheckboxItem
                    value="rich-content"
                    checked={
                      currentDirectoryBlock.config?.show_content ?? false
                    }
                    onCheckedChange={(checked) => {
                      const config = currentDirectoryBlock.config;

                      if (!config) {
                        return;
                      }
                      overwriteBlock({
                        ...currentDirectoryBlock,
                        config: { ...config, show_content: checked },
                      });
                    }}
                  >
                    <HStack w="full">
                      <span>Rich content</span>
                      {currentDirectoryBlock.config?.show_content ? (
                        <ShowIcon />
                      ) : (
                        <HideIcon />
                      )}
                    </HStack>
                  </Menu.CheckboxItem>
                )}
              </Menu.ItemGroup>
            )}

            <Menu.ItemGroup pl="2" py="1">
              <Menu.ItemGroupLabel>Properties</Menu.ItemGroupLabel>

              {columns.map((property) =>
                !property.fixed && supportsCoverImage ? (
                  <PropertyPlacementMenu
                    key={property.fid}
                    property={property}
                  />
                ) : (
                  <Menu.Item key={property.fid} value={property.fid}>
                    {property.hidden ? <HideIcon /> : <ShowIcon />}
                    &nbsp;
                    {property.name}
                  </Menu.Item>
                ),
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function PropertyPlacementMenu({ property }: { property: ColumnDefinition }) {
  const { store } = useLibraryPageContext();
  const block = useDirectoryBlock();
  const positions = [
    ["header", "Header badge"],
    ["details", "Details"],
    ["footer-start", "Bottom-left badge"],
    ["footer-end", "Bottom-right badge"],
  ] as const;
  return (
    <Menu.Root
      lazyMount
      onSelect={({ value }) => {
        const config = block.config;

        if (!config) {
          return;
        }
        store.getState().overwriteBlock({
          ...block,
          config: {
            ...config,
            columns: config.columns.map((column) =>
              column.fid === property.fid
                ? {
                    ...column,
                    hidden: value === "hidden",
                    placement:
                      value === "hidden"
                        ? column.placement
                        : (value as NonNullable<typeof column.placement>),
                  }
                : column,
            ),
          },
        });
      }}
    >
      <Menu.TriggerItem>
        {property.hidden ? <HideIcon /> : <ShowIcon />}&nbsp;{property.name}
      </Menu.TriggerItem>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48">
            <Menu.RadioItemGroup
              value={
                property.hidden ? "hidden" : (property.placement ?? "details")
              }
            >
              {positions.map(([value, label]) => (
                <Menu.RadioItem key={value} value={value}>
                  {label}
                  <Menu.ItemIndicator>
                    <ShowIcon />
                  </Menu.ItemIndicator>
                </Menu.RadioItem>
              ))}
              <Menu.RadioItem value="hidden">
                Hidden
                <Menu.ItemIndicator>
                  <HideIcon />
                </Menu.ItemIndicator>
              </Menu.RadioItem>
            </Menu.RadioItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
