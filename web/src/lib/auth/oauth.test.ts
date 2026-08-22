import { test } from "uvu";
import * as assert from "uvu/assert";

import { formatOAuthGrant } from "./oauth";

test("formats OAuth grant identifiers", () => {
  assert.is(formatOAuthGrant("client_credentials"), "client credentials");
  assert.is(formatOAuthGrant("authorization_code"), "authorization code");
  assert.is(formatOAuthGrant("refresh_token"), "refresh token");
  assert.is(
    formatOAuthGrant("urn:ietf:params:oauth:grant-type:device_code"),
    "device code",
  );
  assert.is(formatOAuthGrant("custom_grant"), "custom_grant");
});

test.run();
