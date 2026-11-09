import assert from "node:assert/strict";
import { test } from "vitest";

import type {
  Account,
  PostReference,
  ThreadReference,
} from "@/api/openapi-schema";
import { Permission } from "@/api/openapi-schema";

import { canDeletePost, canEditPost, canReplyToThread } from "./permissions";

function account(id: string, permissions: string[] = []): Account {
  return {
    id,
    roles: [{ id: "role-1", permissions }],
  } as unknown as Account;
}

function post(authorId: string): PostReference {
  return {
    author: { id: authorId },
  } as unknown as PostReference;
}

test("canDeletePost/canEditPost returns false without account", () => {
  const pr = post("author-1");
  assert.ok(!canDeletePost(pr, undefined));
  assert.ok(!canEditPost(pr, undefined));
});

test("canDeletePost/canEditPost allows the author", () => {
  const pr = post("author-1");
  const acc = account("author-1");
  assert.ok(canDeletePost(pr, acc));
  assert.ok(canEditPost(pr, acc));
});

test("canDeletePost/canEditPost allows MANAGE_POSTS even when not author", () => {
  const pr = post("author-1");
  const acc = account("moderator-1", [Permission.MANAGE_POSTS]);
  assert.ok(canDeletePost(pr, acc));
  assert.ok(canEditPost(pr, acc));
});

test("canDeletePost/canEditPost allows ADMINISTRATOR even when not author", () => {
  const pr = post("author-1");
  const acc = account("admin-1", [Permission.ADMINISTRATOR]);
  assert.ok(canDeletePost(pr, acc));
  assert.ok(canEditPost(pr, acc));
});

test("canDeletePost/canEditPost denies non-author without permission", () => {
  const pr = post("author-1");
  const acc = account("user-1", [Permission.READ_PROFILE]);
  assert.ok(!canDeletePost(pr, acc));
  assert.ok(!canEditPost(pr, acc));
});

test("locked threads only allow replies from Manage Posts or administrators", () => {
  const unlocked = { locked_at: undefined } as ThreadReference;
  const locked = { locked_at: "2026-11-09T10:00:00Z" } as ThreadReference;

  for (const acc of [account("author-1"), account("other-1")]) {
    assert.ok(canReplyToThread(unlocked, acc));
    assert.ok(!canReplyToThread(locked, acc));
  }

  for (const permission of [
    Permission.MANAGE_POSTS,
    Permission.ADMINISTRATOR,
  ]) {
    const acc = account("moderator-1", [permission]);
    assert.ok(canReplyToThread(unlocked, acc));
    assert.ok(canReplyToThread(locked, acc));
  }

  assert.ok(!canReplyToThread(unlocked));
  assert.ok(!canReplyToThread(locked));
});
