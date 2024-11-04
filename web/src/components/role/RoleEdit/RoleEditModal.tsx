import chroma from "chroma-js";
import { uniqueId } from "lodash";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { getRoleListKey, roleCreate } from "@/api/openapi-client/roles";
import { Role } from "@/api/openapi-schema";
import { AddAction } from "@/components/site/Action/Add";
import { EditAction } from "@/components/site/Action/Edit";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { isDefaultRole } from "@/lib/role/defaults";
import { UseDisclosureProps, useDisclosure } from "@/utils/useDisclosure";

import { RoleEditScreen } from "./RoleEditScreen";
import { Props } from "./useRoleEdit";

export function RoleEditModal({
  role,
  onClose,
  onOpen,
  isOpen,
}: UseDisclosureProps & Props) {
  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={`Edit role: ${role.name}`}
    >
      <RoleEditScreen role={role} onSave={onClose} />
    </ModalDrawer>
  );
}

export function RoleEditModalTrigger({ role }: Props) {
  const disclosure = useDisclosure();

  const isDefault = isDefaultRole(role);

  const titleLabel = isDefault ? "You cannot edit a default role" : "Edit role";

  return (
    <>
      <EditAction
        size="xs"
        disabled={isDefault}
        title={titleLabel}
        onClick={disclosure.onOpen}
      />

      <RoleEditModal {...disclosure} role={role} />
    </>
  );
}

export function RoleCreateModalTrigger() {
  const [newRole, setNewRole] = useState<Role | null>(null);
  const disclosure = useDisclosure();
  const { mutate } = useSWRConfig();

  const revalidate = async () => {
    await mutate(getRoleListKey());
  };

  async function handleCreate() {
    await handle(
      async () => {
        const colour = chroma.random().hex();

        const created = await roleCreate({
          name: uniqueId("New role "),
          colour,
          permissions: [],
        });

        revalidate();
        disclosure.onOpen();
        setNewRole(() => created);
      },
      {
        promiseToast: {
          loading: "Creating role...",
          success: "New role created",
        },
      },
    );
  }

  return (
    <>
      <AddAction px="2" variant="subtle" size="xs" onClick={handleCreate}>
        Create role
      </AddAction>

      {newRole && <RoleEditModal {...disclosure} role={newRole} />}
    </>
  );
}
