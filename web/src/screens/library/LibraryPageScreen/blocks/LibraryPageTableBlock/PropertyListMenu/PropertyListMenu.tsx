import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { zipWith } from "lodash";
import { PropsWithChildren } from "react";

import { HideIcon } from "@/components/ui/icons/HideIcon";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import * as Menu from "@/components/ui/menu";
import {
  LibraryPageBlockTypeTable,
  NodeMetadata,
} from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { getDefaultBlockConfig } from "../column";

export function PropertyListMenu({ children }: PropsWithChildren) {
  const { store } = useLibraryPageContext();
  const { setMeta } = store.getState();

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

  async function handleSelect(value: MenuSelectionDetails) {
    const fid = value.value;

    const nextMeta = {
      ...currentMetadata,

      layout: {
        blocks: currentBlocks.map((b) => {
          if (b.type === "table") {
            return {
              ...b,
              config: {
                ...b.config,
                columns:
                  b.config?.columns.map((c) => {
                    if (c.fid === fid) {
                      return {
                        ...c,
                        hidden: !c.hidden,
                      };
                    }
                    return c;
                  }) ?? [],
              },
            };
          }
          return b;
        }),
      },
    } satisfies NodeMetadata;

    setMeta(nextMeta);
  }

  const properties = zipWith(
    currentChildPropertySchema,
    currentTableBlock.config.columns,
    (a, b) => {
      return {
        ...a,
        ...b,
      };
    },
  );

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
          <Menu.Content minW="36">
            <Menu.ItemGroup pl="2" py="1">
              <Menu.ItemGroupLabel>Properties</Menu.ItemGroupLabel>

              {properties.map((property) => (
                <Menu.Item value={property.fid}>
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
