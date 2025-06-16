import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { EyeIcon } from "lucide-react";
import { PropsWithChildren } from "react";

import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";

import { useLibraryPageContext } from "../../Context";

import { ColumnDefinition } from "./column";

type Props = {
  column: ColumnDefinition;
};

export function ColumnMenu({ column, children }: PropsWithChildren<Props>) {
  const { store } = useLibraryPageContext();
  const { setChildPropertyHiddenState, setChildPropertyName } =
    store.getState();

  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "hide-show": {
        handleColumnHide();
        break;
      }
    }
  }

  function handleColumnNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    setChildPropertyName(column.fid, event.target.value);
  }

  function handleColumnHide() {
    setChildPropertyHiddenState(column.fid, !column.hidden);
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
                  // Override Ark.Menu hooking events
                  onKeyDown={(e) => {
                    // Stop arrow keys, space, etc. from bubbling to the menu
                    e.stopPropagation();
                  }}
                  onClick={(e) => {
                    e.stopPropagation();
                  }}
                  onFocus={(e) => {
                    e.stopPropagation();
                  }}
                />
              </Menu.ItemGroup>
            )}

            <Menu.ItemGroup>
              <Menu.Item value="hide-show">
                <EyeIcon />
                &nbsp;Hide column
              </Menu.Item>

              {/* {!column.fixed && (
                <Menu.Item value="delete">
                  <DeleteIcon />
                  &nbsp;Delete
                </Menu.Item>
              )} */}

              {/* TODO: Filtering on child API */}
              {/* <Menu.Item value="filter">
                <FilterIcon />
                &nbsp;Filter...
              </Menu.Item> */}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
