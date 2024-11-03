"use client";

import { MenuSelectionDetails } from "@ark-ui/react";

import { handle } from "@/api/client";
import { Identifier, ProfileReference } from "@/api/openapi-schema";

import { useMemberRoles } from "./useMemberRoles";

export type Props = {
  profile: ProfileReference;
};

export function useMemberRoleMenu(props: Props) {
  const { ready, data, error, addRole, removeRole, revalidate } =
    useMemberRoles(props);
  if (!ready) {
    return {
      ready: false as const,
      error,
    };
  }

  const { roles } = data;

  const handleAdd = (id: Identifier) => {
    handle(
      async () => {
        await addRole(id);
      },
      {
        async cleanup() {
          await revalidate();
        },
      },
    );
  };

  const handleRemove = (id: Identifier) => {
    handle(
      async () => {
        await removeRole(id);
      },
      {
        async cleanup() {
          await revalidate();
        },
      },
    );
  };

  const handleSelect = (details: MenuSelectionDetails) => {
    const id = details.value;
    const role = roles.find((r) => r.id === id);

    if (!role) return;

    if (role.selected) {
      handleRemove(id);
    } else {
      handleAdd(id);
    }
  };

  return {
    ready: true as const,
    data: {
      roles,
    },
    handlers: {
      handleSelect,
    },
  };
}
