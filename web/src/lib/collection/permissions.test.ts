import { test } from "uvu";
import * as assert from "uvu/assert";

import type { Account, Collection } from "@/api/openapi-schema";
import { Permission } from "@/api/openapi-schema";

import { canDeleteCollection, canEditCollection } from "./permissions";

function account(id: string, permissions: string[] = []): Account {
  return {
    id,
    roles: [{ id: "role-1", permissions }],
  } as unknown as Account;
}

function collection(ownerId: string): Collection {
  return {
    owner: { id: ownerId },
  } as unknown as Collection;
}

test("canDeleteCollection/canEditCollection returns false without account", () => {
  const col = collection("owner-1");
  assert.not.ok(canDeleteCollection(col, undefined));
  assert.not.ok(canEditCollection(col, undefined));
});

test("canDeleteCollection/canEditCollection allows the owner", () => {
  const col = collection("owner-1");
  const acc = account("owner-1");
  assert.ok(canDeleteCollection(col, acc));
  assert.ok(canEditCollection(col, acc));
});

test("canDeleteCollection/canEditCollection allows MANAGE_COLLECTIONS", () => {
  const col = collection("owner-1");
  const acc = account("moderator-1", [Permission.MANAGE_COLLECTIONS]);
  assert.ok(canDeleteCollection(col, acc));
  assert.ok(canEditCollection(col, acc));
});

test("canDeleteCollection/canEditCollection denies non-owner without permission", () => {
  const col = collection("owner-1");
  const acc = account("user-1", [Permission.READ_COLLECTION]);
  assert.not.ok(canDeleteCollection(col, acc));
  assert.not.ok(canEditCollection(col, acc));
});

test.run();
