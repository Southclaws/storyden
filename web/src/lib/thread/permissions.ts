import { Account, Permission, PostReference } from "@/api/openapi-schema";
import { hasPermissionOr } from "@/utils/permissions";

export function canDeletePost(pr: PostReference, account?: Account) {
  if (!account) return false;

  return hasPermissionOr(
    account,
    () => pr.author.id === account.id,
    Permission.MANAGE_POSTS,
  );
}

export function canEditPost(pr: PostReference, account?: Account) {
  if (!account) return false;

  return hasPermissionOr(
    account,
    () => pr.author.id === account.id,
    Permission.MANAGE_POSTS,
  );
}
