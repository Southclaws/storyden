import { zodResolver } from "@hookform/resolvers/zod";
import assert from "node:assert/strict";
import { test } from "vitest";

import * as z from "@/lib/validation";

import {
  ExistingPasswordSchema,
  PasswordSchema,
  UsernameOrEmailSchema,
  UsernameSchema,
} from "./schemas";

test("usernames normalize uppercase before validating the handle", () => {
  assert.strictEqual(UsernameSchema.parse("Alice_123"), "alice_123");
  assert.strictEqual(UsernameSchema.parse("a".repeat(30)), "a".repeat(30));
  assert.ok(!UsernameSchema.safeParse("a".repeat(31)).success);
  assert.ok(!UsernameSchema.safeParse("alice smith").success);
  assert.ok(!UsernameSchema.safeParse(" alice ").success);
});

test("password validation preserves its field messages and English type errors", () => {
  const short = PasswordSchema.safeParse("1234567");
  assert.ok(!short.success);
  assert.strictEqual(
    short.error.issues[0]?.message,
    "Password must be at least 8 characters.",
  );
  assert.strictEqual(PasswordSchema.parse("12345678"), "12345678");

  const missing = ExistingPasswordSchema.safeParse("");
  assert.ok(!missing.success);
  assert.strictEqual(
    missing.error.issues[0]?.message,
    "Please enter your current password.",
  );

  const wrongType = PasswordSchema.safeParse(123);
  assert.ok(!wrongType.success);
  assert.strictEqual(
    wrongType.error.issues[0]?.message,
    "Invalid input: expected string, received number",
  );
});

test("sign-in identifiers accept handles and email addresses", () => {
  assert.strictEqual(UsernameOrEmailSchema.parse("ALICE"), "alice");
  assert.strictEqual(
    UsernameOrEmailSchema.parse("Alice@example.com"),
    "Alice@example.com",
  );
  assert.ok(!UsernameOrEmailSchema.safeParse("alice@invalid").success);
});

test("React Hook Form receives normalized values and strips unknown fields", async () => {
  const resolver = zodResolver(
    z.object({ identifier: UsernameSchema, token: PasswordSchema }),
  );
  const input = { identifier: "ALICE", token: "password123", extra: true };
  const result = await resolver(input, undefined, {
    fields: {},
    shouldUseNativeValidation: false,
  });
  assert.deepStrictEqual(result, {
    values: { identifier: "alice", token: "password123" },
    errors: {},
  });
});

test("React Hook Form maps Mini issues to the original fields and messages", async () => {
  const resolver = zodResolver(
    z.object({ identifier: UsernameSchema, token: PasswordSchema }),
  );
  const result = await resolver({ identifier: "", token: "short" }, undefined, {
    fields: {},
    shouldUseNativeValidation: false,
    criteriaMode: "all",
  });
  assert.deepStrictEqual(result.values, {});
  assert.strictEqual(
    result.errors.identifier?.message,
    "Please enter a username.",
  );
  assert.strictEqual(
    result.errors.token?.message,
    "Password must be at least 8 characters.",
  );
  assert.deepStrictEqual(result.errors.identifier?.types, {
    too_small: "Please enter a username.",
  });
});
