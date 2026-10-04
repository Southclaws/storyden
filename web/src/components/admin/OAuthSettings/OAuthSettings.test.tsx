import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import { OAuthSettings } from "./OAuthSettings";

const mocks = vi.hoisted(() => ({ update: vi.fn() }));
vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ updateSettings: mocks.update }),
}));
vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<unknown>) => fn(),
}));

beforeEach(() => {
  vi.clearAllMocks();
  mocks.update.mockResolvedValue(undefined);
});

it("keeps agent management separate from OAuth client settings", () => {
  render(
    <OAuthSettings
      settings={{ ...DefaultSettings, capabilities: ["oauth"] }}
      clients={[]}
      deviceAuthorisations={[]}
      tokens={[]}
      tokenPage={{ currentPage: 1, totalPages: 0, pageSize: 50 }}
    />,
  );
  expect(screen.getByRole("heading", { name: "OAuth" })).toBeVisible();
  expect(screen.getByText("Clients")).toBeVisible();
  expect(screen.queryByText("Agent registrations")).not.toBeInTheDocument();
  expect(
    screen.queryByRole("link", { name: /registration/i }),
  ).not.toBeInTheDocument();
});

it.each([false, true])(
  "saves only the DCR toggle when its saved value is %s",
  async (enabled) => {
    render(
      <OAuthSettings
        settings={{
          ...DefaultSettings,
          capabilities: ["oauth"],
          services: {
            oauth: {
              dynamic_registration_enabled: enabled,
              autonomous_registration_mode: "approval",
              autonomous_registration_role_id: "dasf64to2dtl50ds1n00",
            },
          },
        }}
        clients={[]}
        deviceAuthorisations={[]}
        tokens={[]}
        tokenPage={{ currentPage: 1, totalPages: 0, pageSize: 50 }}
      />,
    );
    const toggle = screen.getByRole("checkbox", {
      name: "Enable dynamic client registration",
    });
    if (enabled) expect(toggle).toBeChecked();
    else expect(toggle).not.toBeChecked();
    await userEvent.click(toggle);
    expect(mocks.update).not.toHaveBeenCalled();
    await userEvent.click(screen.getByRole("button", { name: "Save" }));
    await waitFor(() =>
      expect(mocks.update).toHaveBeenCalledWith({
        services: { oauth: { dynamic_registration_enabled: !enabled } },
      }),
    );
  },
);

it("disables DCR configuration when OAuth is unavailable", async () => {
  render(
    <OAuthSettings
      settings={{ ...DefaultSettings, capabilities: [] }}
      clients={[]}
      deviceAuthorisations={[]}
      tokens={[]}
      tokenPage={{ currentPage: 1, totalPages: 0, pageSize: 50 }}
    />,
  );
  expect(
    screen.getByRole("checkbox", {
      name: "Enable dynamic client registration",
    }),
  ).toBeDisabled();
  const save = screen.getByRole("button", { name: "Save" });
  expect(save).toBeDisabled();
  await userEvent.click(save);
  expect(mocks.update).not.toHaveBeenCalled();
});
