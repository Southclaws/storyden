"use client";

import { TrashIcon } from "@heroicons/react/24/outline";

import { CancelAction } from "@/components/site/Action/Cancel";
import * as Menu from "@/components/ui/menu";
import { HStack } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

export type Props = {
  isConfirmingDelete: boolean;

  onAttemptDelete: () => void;
  onCancelDelete: () => void;
  onDelete?: () => void;
};

export function DeleteWithConfirmationMenuItem(props: Props) {
  if (props.isConfirmingDelete) {
    return (
      <HStack gap="0">
        <Menu.Item
          className={menuItemColorPalette({ colorPalette: "red" })}
          value="delete"
          w="full"
          closeOnSelect={false}
          onClick={props.onDelete}
        >
          Are you sure?
        </Menu.Item>

        <Menu.Item value="delete-cancel" closeOnSelect={false} asChild>
          <CancelAction borderRadius="md" onClick={props.onCancelDelete} />
        </Menu.Item>
      </HStack>
    );
  }

  return (
    <Menu.Item
      className={menuItemColorPalette({ colorPalette: "red" })}
      value="delete"
      closeOnSelect={false}
      onClick={props.onAttemptDelete}
    >
      <HStack gap="1">
        <TrashIcon width="1.4em" /> Delete
      </HStack>
    </Menu.Item>
  );
}
