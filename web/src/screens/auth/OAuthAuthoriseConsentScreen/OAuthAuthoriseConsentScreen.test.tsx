import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { RequestError } from "@/api/common";

import { OAuthAuthoriseConsentScreen } from "./OAuthAuthoriseConsentScreen";

let params = new URLSearchParams();
let hookState: unknown;

const mocks = vi.hoisted(() => ({
  replace: vi.fn(),
  submit: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  usePathname: () => "/oauth/authorize/consent",
  useRouter: () => ({ replace: mocks.replace }),
  useSearchParams: () => params,
}));

vi.mock("@/api/openapi-client/auth", () => ({
  useOAuthAuthoriseConsent: () => hookState,
  oAuthAuthoriseConsentSubmit: mocks.submit,
}));

describe("OAuthAuthoriseConsentScreen", () => {
  beforeEach(() => {
    params = new URLSearchParams();
    hookState = { data: undefined, error: undefined };
    mocks.submit.mockReset();
    mocks.replace.mockReset();
  });

  it("asks for the full application link when the request id is missing", () => {
    render(<OAuthAuthoriseConsentScreen />);

    expect(screen.getByRole("heading", { name: "Missing request" })).toBeVisible();
    expect(screen.getByText("Open the full link from the application and try again.")).toBeVisible();
  });

  it("shows access denied when the OAuth API returns access_denied", () => {
    params = new URLSearchParams({ request_id: "request-123" });
    hookState = {
      data: undefined,
      error: new RequestError("Permission denied.", 400, {
        trace_id: "unknown",
        type: "urn:storyden:problem:oauth:access-denied",
        title: "Permission denied.",
        detail: "OAuth error: access_denied",
      }),
    };

    render(<OAuthAuthoriseConsentScreen />);

    expect(screen.getByText("Access denied")).toBeVisible();
    expect(
      screen.getByText(
        "Your account does not have permission to use third-party applications. Contact an administrator for access.",
      ),
    ).toBeVisible();
  });

  it("shows client, redirect destination, and requested permissions", () => {
    params = new URLSearchParams({ request_id: "request-123" });
    hookState = {
      data: {
        request_id: "request-123",
        client_id: "third-party-client",
        client_name: "Analytics Bot",
        redirect_uri: "https://client.example/callback",
        expires_at: new Date().toISOString(),
        requested_scopes: ["openid", "profile", "CREATE_POST"],
        granted_scopes: ["openid", "profile", "CREATE_POST"],
        inherits_user_permissions: false,
      },
      error: undefined,
    };

    render(<OAuthAuthoriseConsentScreen />);

    expect(screen.getByRole("heading", { name: "Analytics Bot" })).toBeVisible();
    expect(screen.getByText("https://client.example/callback")).toBeVisible();
    expect(screen.getByText("Create post")).toBeVisible();
    expect(screen.getByText("Members can create posts.")).toBeVisible();
  });
});
