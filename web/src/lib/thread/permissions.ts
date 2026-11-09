import {
  Account,
  Permission,
  PostReference,
  ThreadReference,
} from "@/api/openapi-schema";
import { hasPermission, hasPermissionOr } from "@/utils/permissions";

export function canReplyToThread(thread: ThreadReference, account?: Account) {
  if (!account) {
    return false;
  }

  if (thread.locked_at) {
    return hasPermission(account, Permission.MANAGE_POSTS);
  }

  return true;
}

export function canDeletePost(pr: PostReference, account?: Account) {
  if (!account) {
    return false;
  }

  return hasPermissionOr(
    account,
    () => pr.author.id === account.id,
    Permission.MANAGE_POSTS,
  );
}

export function canEditPost(pr: PostReference, account?: Account) {
  if (!account) {
    return false;
  }

  return hasPermissionOr(
    account,
    () => pr.author.id === account.id,
    Permission.MANAGE_POSTS,
  );
}
