import { z } from "zod";

import { Permission } from "@/api/openapi-schema";

export const PermissionSchema = z.nativeEnum(Permission);

export type PermissionDetail = {
  value: Permission;
  name: string;
  description: string;
};

export const PermissionDetails: Record<Permission, PermissionDetail> = {
  [Permission.CREATE_POST]: {
    value: Permission.CREATE_POST,
    name: "Create post",
    description: "Members can create posts.",
  },
  [Permission.READ_PUBLISHED_THREADS]: {
    value: Permission.READ_PUBLISHED_THREADS,
    name: "Read published threads",
    description: "Members can read published threads.",
  },
  [Permission.CREATE_REACTION]: {
    value: Permission.CREATE_REACTION,
    name: "Create reactions",
    description: "React to posts.",
  },
  [Permission.MANAGE_POSTS]: {
    value: Permission.MANAGE_POSTS,
    name: "Manage posts",
    description: "Manage posts, such as delete, pin and move.",
  },
  [Permission.MANAGE_CATEGORIES]: {
    value: Permission.MANAGE_CATEGORIES,
    name: "Manage categories",
    description: "Create and edit categories.",
  },
  [Permission.CREATE_INVITATION]: {
    value: Permission.CREATE_INVITATION,
    name: "Create invitations",
    description: "Create invitations for new members.",
  },
  [Permission.READ_PUBLISHED_LIBRARY]: {
    value: Permission.READ_PUBLISHED_LIBRARY,
    name: "Read published library pages",
    description: "Access published items in the library.",
  },
  [Permission.MANAGE_LIBRARY]: {
    value: Permission.MANAGE_LIBRARY,
    name: "Manage library",
    description:
      "Manage items in the library, such as publish new pages, accept or reject submissions and edit pages.",
  },
  [Permission.SUBMIT_LIBRARY_NODE]: {
    value: Permission.SUBMIT_LIBRARY_NODE,
    name: "Submit to library",
    description: "Submit new pages to the library.",
  },
  [Permission.UPLOAD_ASSET]: {
    value: Permission.UPLOAD_ASSET,
    name: "Upload assets",
    description:
      "Upload images and media to posts, pages and any other areas that accept media.",
  },
  [Permission.MANAGE_EVENTS]: {
    value: Permission.MANAGE_EVENTS,
    name: "Manage events",
    description:
      "Create, edit and remove events as well as change event hosts.",
  },
  [Permission.LIST_PROFILES]: {
    value: Permission.LIST_PROFILES,
    name: "Access member list",
    description: "Access the full list of registered members.",
  },
  [Permission.READ_PROFILE]: {
    value: Permission.READ_PROFILE,
    name: "Access profiles",
    description: "Access any registered member's profiles.",
  },
  [Permission.CREATE_COLLECTION]: {
    value: Permission.CREATE_COLLECTION,
    name: "Create collections",
    description: "Create and manage their own collections.",
  },
  [Permission.LIST_COLLECTIONS]: {
    value: Permission.LIST_COLLECTIONS,
    name: "Access published collections",
    description:
      "Access the full list of published collections from other members.",
  },
  [Permission.READ_COLLECTION]: {
    value: Permission.READ_COLLECTION,
    name: "Read any collection",
    description: "Read any published collection from any other member.",
  },
  [Permission.MANAGE_COLLECTIONS]: {
    value: Permission.MANAGE_COLLECTIONS,
    name: "Manage collections",
    description: "Delete, rename or move collections owned by other members.",
  },
  [Permission.COLLECTION_SUBMIT]: {
    value: Permission.COLLECTION_SUBMIT,
    name: "Submit to collections",
    description: "Submit items for review to other members' collections.",
  },
  [Permission.USE_PERSONAL_ACCESS_KEYS]: {
    value: Permission.USE_PERSONAL_ACCESS_KEYS,
    name: "Use personal access keys",
    description:
      "Use personal access keys to authenticate with the Storyden API and MCP server.",
  },
  [Permission.MANAGE_SETTINGS]: {
    value: Permission.MANAGE_SETTINGS,
    name: "Manage settings",
    description:
      "Manage the administrative settings for the Storyden installation.",
  },
  [Permission.MANAGE_SUSPENSIONS]: {
    value: Permission.MANAGE_SUSPENSIONS,
    name: "Manage suspensions",
    description: "Suspend or reinstate members from the community.",
  },
  [Permission.MANAGE_ROLES]: {
    value: Permission.MANAGE_ROLES,
    name: "Manage roles",
    description:
      "Create, edit and delete roles as well as assign and remove roles of other members.",
  },
  [Permission.MANAGE_REPORTS]: {
    value: Permission.MANAGE_REPORTS,
    name: "Manage reports",
    description:
      "View and manage all submitted reports from community members.",
  },
  [Permission.VIEW_ACCOUNTS]: {
    value: Permission.VIEW_ACCOUNTS,
    name: "View accounts",
    description:
      "View detailed account information including email addresses and verification status for non-administrator accounts.",
  },
  [Permission.USE_ROBOTS]: {
    value: Permission.USE_ROBOTS,
    name: "Use robots",
    description:
      "Use Robots to build automations for managing content, moderation and more.",
  },
  [Permission.MANAGE_ROBOTS]: {
    value: Permission.MANAGE_ROBOTS,
    name: "Manage robots",
    description:
      "Create, edit, and delete Robots and manage Robot configurations.",
  },
  [Permission.ADMINISTRATOR]: {
    value: Permission.ADMINISTRATOR,
    name: "Administrator",
    description: "Full administrative access. Use with caution!",
  },
};

export const PermissionList = Object.values(PermissionDetails);

export function buildPermissionList(
  ...permissionNames: Permission[]
): PermissionDetail[] {
  return permissionNames.map((name) => PermissionDetails[name]).filter(Boolean);
}
