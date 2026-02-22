import { test } from "uvu";
import * as assert from "uvu/assert";

import type { Account, PostReference } from "@/api/openapi-schema";
import { Permission } from "@/api/openapi-schema";

import { canDeletePost, canEditPost } from "./permissions";

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
  assert.not.ok(canDeletePost(pr, undefined));
  assert.not.ok(canEditPost(pr, undefined));
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

test("canDeletePost/canEditPost denies non-author without permission", () => {
  const pr = post("author-1");
  const acc = account("user-1", [Permission.READ_PROFILE]);
  assert.not.ok(canDeletePost(pr, acc));
  assert.not.ok(canEditPost(pr, acc));
});

test.run();
