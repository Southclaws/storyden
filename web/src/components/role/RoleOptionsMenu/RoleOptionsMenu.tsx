import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { PropsWithChildren } from "react";

import { handle } from "@/api/client";
import { Role } from "@/api/openapi-schema";
import { DeleteWithConfirmationMenuItem } from "@/components/site/DeleteConfirmationMenuItem";
import { EditIcon } from "@/components/ui/icons/Edit";
import * as Menu from "@/components/ui/menu";
import { HStack } from "@/styled-system/jsx";

import { RoleBadge } from "../RoleBadge/RoleBadge";

export type Props = {
  role: Role;
};

export function RoleOptionsMenu({ children, role }: PropsWithChildren<Props>) {
  function handleSelect(value: MenuSelectionDetails) {
    switch (value.value) {
      case "":
    }
  }

  async function handleDelete() {
    await handle(async () => {
      //
    });
  }

  return (
    <Menu.Root onSelect={handleSelect}>
      <Menu.Trigger cursor="pointer">{children}</Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel display="flex" gap="2" alignItems="center">
                <RoleBadge role={role} />
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="edit">
                <HStack gap="1">
                  <EditIcon />
                  Edit
                </HStack>
              </Menu.Item>

              <DeleteWithConfirmationMenuItem onDelete={handleDelete} />
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
