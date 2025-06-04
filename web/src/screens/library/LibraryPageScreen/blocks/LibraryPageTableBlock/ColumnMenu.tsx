import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { EyeIcon, FilterIcon } from "lucide-react";
import { PropsWithChildren } from "react";

import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import { useEmitLibraryBlockEvent } from "@/lib/library/events";
import { LibraryPageBlock } from "@/lib/library/metadata";

import { useLibraryPageContext } from "../../Context";

import { ColumnDefinition } from "./column";

type Props = {
  column: ColumnDefinition;
};

export function ColumnMenu({ column, children }: PropsWithChildren<Props>) {
  const { node, form } = useLibraryPageContext();
  const emit = useEmitLibraryBlockEvent();

  const currentMetadata = form.watch("meta", node.meta);

  function handleSelect(value: MenuSelectionDetails) {
    emit("library:add-block", {
      type: value.value as LibraryPageBlock["type"],
    });
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect} size="xs">
      <Menu.Trigger asChild>{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup pl="2" py="2">
              <Input size="sm" value={column.name} />
            </Menu.ItemGroup>

            <Menu.ItemGroup>
              <Menu.Item value="hide-show">
                <EyeIcon />
                &nbsp;Hide
              </Menu.Item>
              <Menu.Item value="delete">
                <DeleteIcon />
                &nbsp;Delete
              </Menu.Item>
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
