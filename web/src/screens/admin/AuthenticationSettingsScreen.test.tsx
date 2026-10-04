import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import { AuthenticationSettingsScreen } from "./AuthenticationSettingsScreen";

let oauthEnabled = true;
vi.mock("@/lib/settings/settings-client", () => ({
  useSettings: () => ({
    ready: true,
    settings: {
      ...DefaultSettings,
      capabilities: oauthEnabled ? ["oauth"] : [],
    },
  }),
}));
vi.mock("@/api/openapi-client/admin", () => ({
  useAdminSettingsGet: () => ({
    data: {
      ...DefaultSettings,
      capabilities: undefined,
      services: {
        oauth: {
          dynamic_registration_enabled: true,
          autonomous_registration_mode: "approval",
        },
      },
    },
  }),
}));
vi.mock("@/api/openapi-client/roles", () => ({
  useRoleList: () => ({ data: { roles: [] } }),
}));
vi.mock("@/api/openapi-client/accounts", () => ({
  useAccountGet: () => ({
    data: { roles: [{ permissions: ["ADMINISTRATOR"] }] },
  }),
}));
vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ updateSettings: vi.fn() }),
}));

describe("AuthenticationSettingsScreen", () => {
  beforeEach(() => {
    oauthEnabled = true;
  });
  it("uses public capabilities when admin settings omit them", () => {
    render(<AuthenticationSettingsScreen />);
    expect(screen.getByRole("radio", { name: "Admin approval" })).toBeEnabled();
    expect(
      screen.getByRole("link", { name: "Manage agent registrations" }),
    ).toHaveAttribute("href", "/admin/authentication/agents");
  });
  it("disables saved OAuth policy when the server is disabled", () => {
    oauthEnabled = false;
    render(<AuthenticationSettingsScreen />);
    expect(
      screen.getByRole("radio", { name: "Admin approval" }),
    ).toBeDisabled();
    expect(
      screen.getByRole("combobox", { name: "Default agent role" }),
    ).toBeDisabled();
  });
});
