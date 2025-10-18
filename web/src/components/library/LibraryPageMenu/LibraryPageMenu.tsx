import { MenuOpenChangeDetails, Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import { CancelAction } from "@/components/site/Action/Cancel";
import { ButtonProps } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";
import { HStack, styled } from "@/styled-system/jsx";
import { menuItemColorPalette } from "@/styled-system/patterns";
import { ReportNodeMenuItem } from "@/components/report/ReportNodeMenuItem";

import { Props, useLibraryPageMenu } from "./useLibraryPageMenu";

export function LibraryPageMenu({
  node,
  onClose,
  ...props
}: Props & ButtonProps) {
  const {
    availableOperations,
    deleteEnabled,
    isChildrenHidden,
    isConfirmingDelete,
    isManager,
    handlers,
  } = useLibraryPageMenu({
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
      onSelect={handlers.handleSelect}
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

              <ReportNodeMenuItem node={node} />

              {isManager && (
                <Menu.Item value="toggle-hide-in-tree">
                  {isChildrenHidden
                    ? "Show children in tree"
                    : "Hide children in tree"}
                </Menu.Item>
              )}

              {deleteEnabled &&
                (isConfirmingDelete ? (
                  <HStack gap="0">
                    <Menu.Item
                      className={menuItemColorPalette()}
                      colorPalette="red"
                      value="delete"
                      w="full"
                      closeOnSelect={false}
                    >
                      Are you sure?
                    </Menu.Item>

                    <Menu.Item
                      value="delete-cancel"
                      closeOnSelect={false}
                      asChild
                    >
                      <CancelAction
                        borderRadius="md"
                        onClick={handlers.handleCancelDelete}
                      />
                    </Menu.Item>
                  </HStack>
                ) : (
                  <Menu.Item
                    className={menuItemColorPalette({ colorPalette: "red" })}
                    colorPalette="red"
                    value="delete"
                    closeOnSelect={false}
                  >
                    Delete
                  </Menu.Item>
                ))}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
