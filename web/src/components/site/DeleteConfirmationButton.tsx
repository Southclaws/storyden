"use client";

import { TrashIcon } from "@heroicons/react/24/outline";
import { XIcon } from "lucide-react";

import { Button, ButtonProps } from "@/components/ui/button";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";
import { button } from "@/styled-system/recipes";

import { IconButton } from "../ui/icon-button";

import { useConfirmation } from "./useConfirmation";

export type Props = ButtonProps & {
  onDelete: () => Promise<void>;
};

export function DeleteWithConfirmationButton({ onDelete, ...props }: Props) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onDelete);

  if (isConfirming) {
    return (
      <HStack
        className={cx(
          button(props),
          menuItemColorPalette({ colorPalette: "red" }),
        )}
        px="0"
        w="full"
        gap="0"
      >
        <Button
          type="button"
          className={menuItemColorPalette({ colorPalette: "red" })}
          pl="20"
          w="full"
          title="Confirm delete"
          onClick={onDelete}
        >
          Are you sure?
        </Button>

        <IconButton
          type="button"
          variant="ghost"
          title="Cancel delete"
          onClick={handleCancelAction}
        >
          <XIcon />
        </IconButton>
      </HStack>
    );
  }

  return (
    <Button
      {...props}
      type="button"
      className={menuItemColorPalette({ colorPalette: "red" })}
      title="Delete"
      onClick={handleConfirmAction}
    >
      <HStack gap="1">
        <TrashIcon width="1.4em" /> Delete
      </HStack>
    </Button>
  );
}
