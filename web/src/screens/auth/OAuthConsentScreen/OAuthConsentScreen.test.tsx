import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { OAuthConsentScreen } from "./OAuthConsentScreen";

let params = new URLSearchParams();
let hookState: unknown;

const mocks = vi.hoisted(() => ({
  replace: vi.fn(),
  submit: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  usePathname: () => "/oauth/consent",
  useRouter: () => ({ replace: mocks.replace }),
  useSearchParams: () => params,
}));

vi.mock("@/api/openapi-client/auth", () => ({
  useOAuthDeviceConsent: () => hookState,
  oAuthDeviceConsentSubmit: mocks.submit,
}));

describe("OAuthConsentScreen", () => {
  beforeEach(() => {
    params = new URLSearchParams();
    hookState = { data: undefined, error: undefined, isLoading: false };
    mocks.submit.mockReset();
    mocks.replace.mockReset();
  });

  it("asks for the full application link when the user code is missing", () => {
    render(<OAuthConsentScreen />);

    expect(screen.getByRole("heading", { name: "Missing code" })).toBeVisible();
    expect(screen.getByText("Open the full link from the application and try again.")).toBeVisible();
  });

  it("submits approval and shows the completion state", async () => {
    params = new URLSearchParams({ user_code: "ABCD-EFGH" });
    hookState = {
      data: {
        user_code: "ABCD-EFGH",
        client_id: "storyden-cli",
        client_name: "Storyden",
        expires_at: new Date().toISOString(),
        requested_scopes: ["openid", "profile", "offline_access"],
        granted_scopes: ["openid", "profile", "offline_access", "CREATE_POST"],
        inherits_user_permissions: true,
      },
      error: undefined,
      isLoading: false,
    };
    mocks.submit.mockResolvedValue({ status: "approved" });

    render(<OAuthConsentScreen />);

    expect(screen.getByRole("heading", { name: "Storyden" })).toBeVisible();
    expect(screen.getByText("ABCD-EFGH")).toBeVisible();
    expect(
      screen.getByText("Only approve if this code matches the code shown where you started authentication."),
    ).toBeVisible();
    expect(screen.getByText("Create post")).toBeVisible();
    expect(screen.getByText("Members can create posts.")).toBeVisible();

    await userEvent.click(screen.getByRole("button", { name: "Approve" }));

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Approved" })).toBeVisible();
    });
    expect(screen.getByText("You can return to the application.")).toBeVisible();
    expect(mocks.submit).toHaveBeenCalledWith({
      user_code: "ABCD-EFGH",
      decision: "approve",
    });
  });
});
