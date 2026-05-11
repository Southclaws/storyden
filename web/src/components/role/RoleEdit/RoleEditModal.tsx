import chroma from "chroma-js";
import { uniqueId } from "lodash";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { getRoleListKey, roleCreate } from "@/api/openapi-client/roles";
import { Role } from "@/api/openapi-schema";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { IconButton } from "@/components/ui/icon-button";
import { CreateIcon } from "@/components/ui/icons/Create";
import { EditIcon } from "@/components/ui/icons/Edit";
import { useI18n } from "@/i18n/provider";
import { isDefaultRole, isEditableDefaultRole } from "@/lib/role/defaults";
import { DefaultRoleMetadata } from "@/lib/role/metadata";
import { UseDisclosureProps, useDisclosure } from "@/utils/useDisclosure";

import { RoleEditScreen } from "./RoleEditScreen";
import { Props } from "./useRoleEdit";

export function RoleEditModal({
  role,
  onClose,
  onOpen,
  isOpen,
}: UseDisclosureProps & Props) {
  const { t } = useI18n();

  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={t("Edit role: {{name}}", { name: t(role.name) })}
    >
      <RoleEditScreen role={role} onSave={onClose} />
    </ModalDrawer>
  );
}

export function RoleEditModalTrigger({ role }: Props) {
  const { t } = useI18n();
  const disclosure = useDisclosure();

  const isDefault = isDefaultRole(role);
  const isEditable = isEditableDefaultRole(role);
  const cannotEdit = isDefault && !isEditable;

  const titleLabel = cannotEdit
    ? t("You cannot edit this role")
    : t("Edit role");

  return (
    <>
      <IconButton
        variant="ghost"
        size="xs"
        minWidth="5"
        width="5"
        height="5"
        padding="0"
        color="fg.muted"
        disabled={cannotEdit}
        title={titleLabel}
        onClick={disclosure.onOpen}
      >
        <EditIcon width="4" />
      </IconButton>

      <RoleEditModal {...disclosure} role={role} />
    </>
  );
}

export function RoleCreateModalTrigger() {
  const { t } = useI18n();
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
          meta: DefaultRoleMetadata,
          permissions: [],
        });

        revalidate();
        disclosure.onOpen();
        setNewRole(() => created);
      },
      {
        promiseToast: {
          loading: t("Creating role..."),
          success: t("New role created"),
        },
      },
    );
  }

  return (
    <>
      <IconButton px="2" variant="subtle" size="xs" onClick={handleCreate}>
        <CreateIcon /> {t("Create role")}
      </IconButton>

      {newRole && <RoleEditModal {...disclosure} role={newRole} />}
    </>
  );
}
