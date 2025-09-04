"use client";

import { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";
import { button } from "@/styled-system/recipes";

import { useConfirmation } from "./useConfirmation";

export type Props = ButtonProps & {
  onDelete: () => Promise<void>;
};

export function DeleteWithConfirmationButton({
  onDelete,
  children,
  ...props
}: PropsWithChildren<Props>) {
  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(onDelete);

  if (isConfirming) {
    return (
      <HStack
        className={cx(
          button(props),
          menuItemColorPalette({ colorPalette: "tomato" }),
        )}
        px="0"
        w="full"
        gap="0"
      >
        <Button
          type="button"
          className={menuItemColorPalette({ colorPalette: "tomato" })}
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
          <CancelIcon />
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
        <DeleteIcon /> {children ?? "Delete"}
      </HStack>
    </Button>
  );
}
