import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import type { Account } from "@/api/openapi-schema";
import { parseMemberSettings } from "@/lib/settings/member-settings";

import { MemberInterfaceSettings } from "./MemberInterfaceSettings";

const update = vi.fn();
const revalidate = vi.fn();
const revalidateAccount = vi.fn();

vi.mock("@/api/client", () => ({
  handle: async (
    action: () => Promise<void>,
    options?: { cleanup?: () => Promise<void> },
  ) => {
    await action();
    await options?.cleanup?.();
  },
}));

vi.mock("@/api/openapi-client/accounts", () => ({
  useAccountGet: () => ({ mutate: revalidateAccount }),
}));

vi.mock("@/lib/profile/mutation", () => ({
  useProfileMutations: () => ({ update, revalidate }),
}));

const account: Account = {
  id: "member-1",
  createdAt: "2026-11-07T12:00:00Z",
  updatedAt: "2026-11-07T12:00:00Z",
  joined: "2026-11-07T12:00:00Z",
  handle: "member",
  name: "Member",
  roles: [],
  bio: "",
  links: [],
  meta: { editor: { mode: "markdown" }, custom: "keep" },
  verified_status: "none",
  kind: "human",
  email_addresses: [],
  auth_services: [],
  admin: false,
};

describe("MemberInterfaceSettings", () => {
  beforeEach(() => {
    update.mockResolvedValue(undefined);
    revalidate.mockResolvedValue(undefined);
    revalidateAccount.mockResolvedValue(undefined);
  });

  it("saves the preference with existing settings and refreshes the session", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings(account);

    render(<MemberInterfaceSettings session={session} />);

    const checkbox = screen.getByRole("checkbox", {
      name: "Mark notifications as read when opened",
    });

    expect(checkbox).toBeChecked();

    await user.click(checkbox);

    expect(update).not.toHaveBeenCalled();

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith({
        meta: {
          editor: { mode: "markdown" },
          notifications: { markReadOnClick: false },
          custom: "keep",
        },
      });
      expect(revalidateAccount).toHaveBeenCalledOnce();
      expect(revalidate).toHaveBeenCalledOnce();
    });
  });

  it("loads a saved disabled preference and allows enabling it again", async () => {
    const user = userEvent.setup();
    const session = parseMemberSettings({
      ...account,
      meta: { ...account.meta, notifications: { markReadOnClick: false } },
    });

    render(<MemberInterfaceSettings session={session} />);

    const checkbox = screen.getByRole("checkbox", {
      name: "Mark notifications as read when opened",
    });

    expect(checkbox).not.toBeChecked();

    await user.click(checkbox);
    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(update).toHaveBeenCalledWith({
        meta: {
          editor: { mode: "markdown" },
          notifications: { markReadOnClick: true },
          custom: "keep",
        },
      });
    });
  });
});
