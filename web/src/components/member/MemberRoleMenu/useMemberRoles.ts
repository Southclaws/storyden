"use client";

import { keyBy } from "lodash";

import {
  accountAddRole,
  accountRemoveRole,
} from "@/api/openapi-client/accounts";
import { useProfileGet } from "@/api/openapi-client/profiles";
import { useRoleList } from "@/api/openapi-client/roles";
import {
  AccountRole,
  Identifier,
  ProfileReference,
} from "@/api/openapi-schema";

export type Props = {
  profile: ProfileReference;
};

export function useMemberRoles(props: Props) {
  const accountID = props.profile.id;
  const accountHandle = props.profile.handle;

  const { data: roleData, error: roleError } = useRoleList();

  const {
    data: profileData,
    error: profileError,
    mutate,
  } = useProfileGet(accountID);

  if (!roleData || !profileData) {
    return {
      ready: false as const,
      error: roleError || profileError,
    };
  }

  const currentRoles = keyBy(profileData.roles, "id");

  const roles = roleData.roles.map((r) => {
    const selected = !!currentRoles[r.id];
    return {
      ...r,
      selected,
    };
  });

  const addRole = async (id: Identifier) => {
    await mutate(
      (prev) => {
        const current = prev ?? profileData;

        const newRole = roleData.roles.find((r) => r.id === id);
        if (!newRole) return;

        const newHeldRole = { ...newRole } as AccountRole;
        const nextRoles = [...current.roles, newHeldRole];

        const nextProfile = {
          ...current,
          roles: nextRoles,
        };

        return nextProfile;
      },
      { revalidate: false },
    );

    await accountAddRole(accountHandle, id);
  };

  const removeRole = async (id: Identifier) => {
    await mutate(
      (prev) => {
        const current = prev ?? profileData;

        const nextRoles = [...current.roles.filter((r) => r.id != id)];

        const nextProfile = {
          ...current,
          roles: nextRoles,
        };

        return nextProfile;
      },
      { revalidate: false },
    );

    await accountRemoveRole(accountHandle, id);
  };

  return {
    ready: true as const,
    data: {
      roles,
    },
    addRole,
    removeRole,
    revalidate: mutate,
  };
}
