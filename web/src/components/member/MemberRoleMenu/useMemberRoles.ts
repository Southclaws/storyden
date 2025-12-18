"use client";

import { keyBy } from "lodash";
import { useState } from "react";

import {
  accountAddRole,
  accountRemoveRole,
} from "@/api/openapi-client/accounts";
import { useProfileGet } from "@/api/openapi-client/profiles";
import { useRoleList } from "@/api/openapi-client/roles";
import { Identifier, ProfileReference } from "@/api/openapi-schema";

export type Props = {
  profile: ProfileReference;
};

export function useMemberRoles(props: Props) {
  const accountID = props.profile.id;
  const accountHandle = props.profile.handle;

  const [isUpdating, setIsUpdating] = useState(false);

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
    setIsUpdating(true);
    try {
      await accountAddRole(accountHandle, id);
      await mutate();
    } finally {
      setIsUpdating(false);
    }
  };

  const removeRole = async (id: Identifier) => {
    setIsUpdating(true);
    try {
      await accountRemoveRole(accountHandle, id);
      await mutate();
    } finally {
      setIsUpdating(false);
    }
  };

  return {
    ready: true as const,
    data: {
      roles,
    },
    addRole,
    removeRole,
    revalidate: mutate,
    isUpdating,
  };
}
