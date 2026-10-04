import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import {
  AuthenticationSettingsForm,
  type Props,
} from "./AuthenticationSettingsForm";

const mocks = vi.hoisted(() => ({ update: vi.fn() }));
vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ updateSettings: mocks.update }),
}));
vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<unknown>) => fn(),
}));

const roleID = "dasf64to2dtl50ds1n00";
function props(): Props {
  return {
    settings: {
      ...DefaultSettings,
      capabilities: ["oauth"],
      services: {
        oauth: {
          dynamic_registration_enabled: true,
          autonomous_registration_mode: "approval",
          autonomous_registration_role_id: roleID,
        },
      },
    },
    roles: [
      {
        id: roleID,
        name: "Bot",
        colour: "#123456",
        permissions: ["MANAGE_LIBRARY"],
        createdAt: "2026-09-27T12:00:00Z",
        updatedAt: "2026-09-27T12:00:00Z",
      },
    ],
    canManageRoles: true,
    canApproveRegistrations: true,
  };
}
async function save() {
  await userEvent.click(screen.getAllByRole("button", { name: "Save" })[0]!);
  await waitFor(() => expect(mocks.update).toHaveBeenCalled());
}

describe("AuthenticationSettingsForm", () => {
  beforeEach(() => {
    mocks.update.mockResolvedValue(undefined);
  });

  it("disables OAuth fields when unavailable and leaves their saved policy untouched", async () => {
    const input = props();
    input.settings.capabilities = [];
    render(<AuthenticationSettingsForm {...input} />);
    expect(
      screen.queryByRole("checkbox", {
        name: "Enable dynamic client registration",
      }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByRole("combobox", { name: "Default agent role" }),
    ).toBeDisabled();
    expect(
      screen.getByRole("radio", { name: "Admin approval" }),
    ).toBeDisabled();
    await save();
    expect(mocks.update).toHaveBeenCalledWith({
      authentication_mode: "handle",
      registration_mode: "public",
    });
  });

  it("saves autonomous admission separately from ordinary client registration", async () => {
    render(<AuthenticationSettingsForm {...props()} />);
    await userEvent.click(
      screen.getByRole("radio", { name: "Registration token required" }),
    );
    await save();
    expect(mocks.update).toHaveBeenCalledWith({
      authentication_mode: "handle",
      registration_mode: "public",
      services: {
        oauth: {
          autonomous_registration_mode: "protected",
        },
      },
    });
  });

  it("reads DCR availability from settings without editing it here", () => {
    const input = props();
    input.settings.services!.oauth!.dynamic_registration_enabled = false;
    const { rerender } = render(<AuthenticationSettingsForm {...input} />);
    expect(
      screen.getByRole("radio", { name: "Admin approval" }),
    ).toBeDisabled();
    expect(
      screen.queryByRole("checkbox", {
        name: "Enable dynamic client registration",
      }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByText(/Enable dynamic client registration in OAuth settings/),
    ).toBeVisible();
    rerender(<AuthenticationSettingsForm {...props()} />);
    expect(screen.getByRole("radio", { name: "Admin approval" })).toBeEnabled();
  });

  it("clears the default role explicitly with null", async () => {
    render(<AuthenticationSettingsForm {...props()} />);
    const select = document.querySelector(
      'select[name="autonomous_registration_role_id"]',
    )!;
    fireEvent.change(select, { target: { value: "" } });
    await save();
    expect(
      mocks.update.mock.calls[0]![0].services.oauth
        .autonomous_registration_role_id,
    ).toBeNull();
  });

  it("preserves a deleted role until the administrator changes it", async () => {
    const input = props();
    input.roles = [];
    render(<AuthenticationSettingsForm {...input} />);
    expect(screen.getByText(/The saved role was deleted/)).toBeVisible();
    expect(
      document.querySelector('select[name="autonomous_registration_role_id"]'),
    ).toHaveValue(roleID);
    await save();
    expect(mocks.update.mock.calls[0]![0].services.oauth).not.toHaveProperty(
      "autonomous_registration_role_id",
    );
  });

  it("lets settings managers save policy without assigning roles", async () => {
    render(
      <AuthenticationSettingsForm
        {...props()}
        canManageRoles={false}
        canApproveRegistrations={false}
      />,
    );
    expect(
      screen.getByRole("combobox", { name: "Default agent role" }),
    ).toBeDisabled();
    expect(
      screen.queryByRole("link", { name: "Manage agent registrations" }),
    ).not.toBeInTheDocument();
    await save();
    expect(mocks.update.mock.calls[0]![0].services.oauth).not.toHaveProperty(
      "autonomous_registration_role_id",
    );
  });
});
