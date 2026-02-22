import { Identifier, Permission, Role } from "@/api/openapi-schema";

const DefaultRoleGuestID = "0000000000000000000g";
const DefaultRoleMemberID = "000000000000000000m0";
const DefaultRoleAdminID = "00000000000000000a00";

export function isDefaultRole(role: { id: Identifier }) {
  switch (role.id) {
    case DefaultRoleGuestID:
    case DefaultRoleMemberID:
    case DefaultRoleAdminID:
      return true;

    default:
      return false;
  }
}

export function isGuestRole(role: { id: Identifier }) {
  return role.id === DefaultRoleGuestID;
}

export function isMemberRole(role: { id: Identifier }) {
  return role.id === DefaultRoleMemberID;
}

// Tells you if a default role has been edited with custom permissions.
export function isStoredDefaultRole(role: Role) {
  if (!isDefaultRole(role)) {
    return false;
  }

  // NOTE: Massive hack because we don't expose when a default role is edited.
  return role.createdAt !== "0001-01-01T00:00:00Z";
}

export function isEditableDefaultRole(role: Role) {
  switch (role.id) {
    case DefaultRoleMemberID:
    case DefaultRoleGuestID:
      return true;

    default:
      return false;
  }
}

export const readPermissions = [
  Permission.READ_PUBLISHED_THREADS,
  Permission.READ_PUBLISHED_LIBRARY,
  Permission.LIST_PROFILES,
  Permission.READ_PROFILE,
  Permission.LIST_COLLECTIONS,
  Permission.READ_COLLECTION,
];

export const writePermissions = [
  Permission.CREATE_POST,
  Permission.CREATE_REACTION,
  Permission.MANAGE_POSTS,
  Permission.MANAGE_CATEGORIES,
  Permission.CREATE_INVITATION,
  Permission.MANAGE_LIBRARY,
  Permission.SUBMIT_LIBRARY_NODE,
  Permission.UPLOAD_ASSET,
  Permission.MANAGE_EVENTS,
  Permission.CREATE_COLLECTION,
  Permission.MANAGE_COLLECTIONS,
  Permission.COLLECTION_SUBMIT,
  Permission.USE_PERSONAL_ACCESS_KEYS,
  Permission.MANAGE_SETTINGS,
  Permission.MANAGE_SUSPENSIONS,
  Permission.MANAGE_ROLES,
  Permission.ADMINISTRATOR,
];

export function isWritePermission(permission: Permission) {
  return writePermissions.includes(
    permission as (typeof writePermissions)[number],
  );
}
export function isReadPermission(permission: Permission) {
  return readPermissions.includes(
    permission as (typeof readPermissions)[number],
  );
}
