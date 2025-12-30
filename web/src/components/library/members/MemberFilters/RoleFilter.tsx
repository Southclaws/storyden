"use client";

import { filter, map } from "lodash/fp";
import { parseAsArrayOf, parseAsString, useQueryState } from "nuqs";
import { useEffect, useMemo, useState } from "react";

import { useRoleList } from "@/api/openapi-client/roles";
import { Role } from "@/api/openapi-schema";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/MultiSelectPicker";
import { isGuestRole, isMemberRole } from "@/lib/role/defaults";

const roleToMultiSelectItem = map(
  (role: Role): MultiSelectPickerItem => ({
    label: role.name,
    value: role.id,
    colour: role.colour,
  }),
);

const filterOutGuestRole = filter<Role>(
  (role: Role) => !isGuestRole(role) && !isMemberRole(role),
);

export function RoleFilter() {
  const [roles, setRoles] = useQueryState(
    "roles",
    parseAsArrayOf(parseAsString).withDefault([]),
  );

  const { data, error } = useRoleList();
  const allRoles = filterOutGuestRole(data?.roles || []);
  const items = roleToMultiSelectItem(allRoles);

  const [searchResults, setSearchResults] = useState<MultiSelectPickerItem[]>(
    [],
  );
  const [hasUserQueried, setHasUserQueried] = useState(false);

  useEffect(() => {
    if (items.length > 0 && !hasUserQueried) {
      setSearchResults(items);
    }
  }, [items, hasUserQueried]);

  const selectedRoles = useMemo(() => {
    return roles
      .map((roleId) => items.find((r) => r.value === roleId))
      .filter((r): r is MultiSelectPickerItem => r !== undefined);
  }, [roles, items]);

  function handleQuery(query: string) {
    setHasUserQueried(true);

    if (query.length === 0) {
      setSearchResults(items);
      return;
    }

    const lowerQuery = query.toLowerCase();
    const filtered = items.filter((role) =>
      role.label.toLowerCase().includes(lowerQuery),
    );
    setSearchResults(filtered);
  }

  async function handleChange(items: MultiSelectPickerItem[]) {
    const roleIds = items.map((item) => item.value);
    await setRoles(roleIds.length > 0 ? roleIds : null);
  }

  return (
    <MultiSelectPicker
      inputPlaceholder="Filter roles"
      value={selectedRoles}
      onChange={handleChange}
      onQuery={handleQuery}
      queryResults={searchResults}
      queryError={error?.message}
      size="sm"
      triggerProps={{
        width: "full",
        minW: "32",
        flexShrink: "1",
      }}
    />
  );
}
