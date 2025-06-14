import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { HideIcon } from "@/components/ui/icons/HideIcon";
import { ShowIcon } from "@/components/ui/icons/ShowIcon";
import * as Menu from "@/components/ui/menu";

import { useLibraryPageContext } from "../../../Context";
import { useWatch } from "../../../store";
import { mergeFieldsAndPropertySchema } from "../column";
import { useTableBlock } from "../useTableBlock";

export function PropertyListMenu({ children }: PropsWithChildren) {
  const { store } = useLibraryPageContext();
  const { setChildPropertyHiddenState } = store.getState();

  const currentTableBlock = useTableBlock();
  const currentChildPropertySchema = useWatch(
    (s) => s.draft.child_property_schema,
  );

  function handleSelect(value: MenuSelectionDetails) {
    const fid = value.value;

    const hidden =
      currentTableBlock.config?.columns?.find((c) => c.fid === fid)?.hidden ??
      true;

    setChildPropertyHiddenState(fid, !hidden);
  }

  const columns = mergeFieldsAndPropertySchema(
    currentChildPropertySchema,
    currentTableBlock,
    true,
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
