"use client";

import { TrashIcon } from "@heroicons/react/24/outline";

import { CancelAction } from "@/components/site/Action/Cancel";
import * as Menu from "@/components/ui/menu";
import { HStack } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";

import { useConfirmation } from "./useConfirmation";

export type Props = {
  onDelete: () => Promise<void>;
};

export function DeleteWithConfirmationMenuItem(props: Props) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(props.onDelete);

  if (isConfirming) {
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
          <CancelAction borderRadius="md" onClick={handleCancelAction} />
        </Menu.Item>
      </HStack>
    );
  }

  return (
    <Menu.Item
      className={menuItemColorPalette({ colorPalette: "red" })}
      value="delete"
      closeOnSelect={false}
      onClick={handleConfirmAction}
    >
      <HStack gap="1">
        <TrashIcon width="1.4em" /> Delete
      </HStack>
    </Menu.Item>
  );
}
