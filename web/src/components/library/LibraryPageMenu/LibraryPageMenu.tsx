import { MenuOpenChangeDetails, Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { DeleteConfirmation } from "src/components/site/Action/Delete";
import { MoreAction } from "src/components/site/Action/More";

import { ButtonProps } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";
import { styled } from "@/styled-system/jsx";

import { Props, useLibraryPageMenu } from "./useLibraryPageMenu";

export function LibraryPageMenu({
  node,
  onClose,
  ...props
}: Props & ButtonProps) {
  const { availableOperations, deleteEnabled, deleteProps, handleSelect } =
    useLibraryPageMenu({
      node,
      onClose,
    });

  function handleOpenChange(d: MenuOpenChangeDetails) {
    if (!d.open) {
      onClose?.();
    }
  }

  const statusText =
    node.visibility === "draft"
      ? "(draft)"
      : node.visibility === "review"
        ? "(in review)"
        : "";

  return (
    <Menu.Root
      lazyMount
      positioning={{ placement: "right-start", gutter: -2 }}
      onSelect={handleSelect}
      onOpenChange={handleOpenChange}
    >
      <Menu.Trigger asChild>
        <MoreAction variant="subtle" size="xs" {...props} />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup>
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>
                  {`Created by ${node.owner.name}`} {statusText}
                </styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(node.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              {availableOperations.map((op) => (
                <Menu.Item
                  key={op.targetVisibility}
                  value={op.targetVisibility}
                >
                  {op.label}
                </Menu.Item>
              ))}

              {deleteEnabled && (
                <>
                  <Menu.Item value="delete">Delete</Menu.Item>
                  <DeleteConfirmation {...deleteProps} />
                </>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
