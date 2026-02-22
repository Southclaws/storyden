import { test } from "uvu";
import * as assert from "uvu/assert";

import type { Account } from "@/api/openapi-schema";
import { Permission } from "@/api/openapi-schema";

import { hasPermission, hasPermissionOr } from "./permissions";

function accountWithPermissions(...permissions: string[]): Account {
  return {
    id: "acc_1",
    roles: [
      {
        id: "role_1",
        permissions,
      },
    ],
  } as unknown as Account;
}

test("hasPermission returns false without account", () => {
  assert.not.ok(hasPermission(undefined, Permission.MANAGE_POSTS));
});

test("hasPermission checks permissions across roles", () => {
  const account = {
    id: "acc_1",
    roles: [
      { id: "role_1", permissions: [Permission.READ_PROFILE] },
      { id: "role_2", permissions: [Permission.MANAGE_POSTS] },
    ],
  } as unknown as Account;

  assert.ok(hasPermission(account, Permission.MANAGE_POSTS));
  assert.not.ok(hasPermission(account, Permission.MANAGE_SETTINGS));
});

test("hasPermission grants all when account has ADMINISTRATOR", () => {
  const account = accountWithPermissions(Permission.ADMINISTRATOR);
  assert.ok(hasPermission(account, Permission.MANAGE_SETTINGS));
  assert.ok(hasPermission(account, Permission.MANAGE_POSTS));
});

test("hasPermissionOr returns true from fallback when no permission", () => {
  const account = accountWithPermissions(Permission.READ_PROFILE);
  const result = hasPermissionOr(
    account,
    () => true,
    Permission.MANAGE_POSTS,
  );
  assert.ok(result);
});

test("hasPermissionOr returns false without account even if fallback true", () => {
  assert.not.ok(
    hasPermissionOr(undefined, () => true, Permission.MANAGE_POSTS),
  );
});

test.run();
