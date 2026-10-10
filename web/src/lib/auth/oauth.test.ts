import assert from "node:assert/strict";
import { test } from "vitest";

import { formatOAuthGrant } from "./oauth";

test("formats OAuth grant identifiers", () => {
  assert.strictEqual(
    formatOAuthGrant("client_credentials"),
    "client credentials",
  );
  assert.strictEqual(
    formatOAuthGrant("authorization_code"),
    "authorization code",
  );
  assert.strictEqual(formatOAuthGrant("refresh_token"), "refresh token");
  assert.strictEqual(
    formatOAuthGrant("urn:ietf:params:oauth:grant-type:device_code"),
    "device code",
  );
  assert.strictEqual(formatOAuthGrant("custom_grant"), "custom_grant");
});
