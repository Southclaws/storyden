import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { toast } from "sonner";
import useSWR, { SWRConfig } from "swr";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { getThreadListKey } from "@/api/openapi-client/threads";
import type {
  Account,
  ThreadListOKResponse,
  ThreadReference,
} from "@/api/openapi-schema";

import { LikeButton } from "./LikeButton";

const api = vi.hoisted(() => ({ add: vi.fn(), remove: vi.fn() }));

vi.mock("@/api/openapi-client/likes", () => ({
  likePostAdd: api.add,
  likePostRemove: api.remove,
}));
vi.mock("@/auth", () => ({ useSession: (initial?: Account) => initial }));
vi.mock("sonner", () => ({ toast: { error: vi.fn() } }));

const account = { id: "reader" } as Account;
const thread = {
  id: "harbour",
  title: "Harbour notes",
  likes: { likes: 3, liked: false },
} as ThreadReference;
let serverPage: ThreadListOKResponse;

function Feed({ signedIn = true }: { signedIn?: boolean }) {
  const [initialPage] = useState(serverPage);
  const { data } = useSWR(
    getThreadListKey({ page: "1" }),
    async () => serverPage,
    { fallbackData: initialPage },
  );
  return (
    <LikeButton
      thread={data!.threads[0]!}
      initialSession={signedIn ? account : undefined}
      label="Like Harbour notes"
      showCount
    />
  );
}

function renderFeed(signedIn = true) {
  return render(<Feed signedIn={signedIn} />, {
    wrapper: ({ children }) => (
      <SWRConfig
        value={{
          provider: () => new Map(),
          revalidateOnMount: false,
          revalidateOnFocus: false,
          dedupingInterval: 0,
        }}
      >
        {children}
      </SWRConfig>
    ),
  });
}

describe("feed likes", () => {
  beforeEach(() => {
    serverPage = { threads: [thread] } as ThreadListOKResponse;
  });

  it("reports likes and the sign-in requirement without an unusable guest button", () => {
    renderFeed(false);
    expect(screen.queryByRole("button")).not.toBeInTheDocument();
    expect(
      screen.getByText(/3 likes\. Sign in to like discussions\./),
    ).toBeInTheDocument();
  });

  it("keeps the action name and focus while liking and undoing a like", async () => {
    const user = userEvent.setup();
    api.add.mockImplementation(async () => {
      serverPage = {
        threads: [{ ...thread, likes: { likes: 4, liked: true } }],
      } as ThreadListOKResponse;
    });
    api.remove.mockImplementation(async () => {
      serverPage = { threads: [thread] } as ThreadListOKResponse;
    });
    renderFeed();
    const button = screen.getByRole("button", {
      name: "Like Harbour notes",
      pressed: false,
    });
    expect(button).toHaveAccessibleDescription("3 likes.");
    await user.tab();
    await user.keyboard("{Enter}");
    await waitFor(() => expect(button).toHaveAttribute("aria-pressed", "true"));
    expect(button).toHaveAccessibleDescription("4 likes.");
    expect(button).toHaveAccessibleName("Like Harbour notes");
    expect(button).toHaveFocus();
    await waitFor(() => expect(button).not.toHaveAttribute("aria-busy"));
    await user.keyboard(" ");
    await waitFor(() =>
      expect(button).toHaveAttribute("aria-pressed", "false"),
    );
    expect(button).toHaveAccessibleDescription("3 likes.");
    expect(button).toHaveFocus();
    expect(api.add).toHaveBeenCalledExactlyOnceWith("harbour");
    expect(api.remove).toHaveBeenCalledExactlyOnceWith("harbour");
  });

  it("retains focus during latency and restores server state after an error", async () => {
    let reject!: (error: Error) => void;
    api.add.mockImplementation(
      () =>
        new Promise((_, rejectRequest) => {
          reject = rejectRequest;
        }),
    );
    const user = userEvent.setup();
    renderFeed();
    const button = screen.getByRole("button", { name: "Like Harbour notes" });
    await user.tab();
    await user.keyboard("{Enter}");
    await waitFor(() => expect(button).toHaveAttribute("aria-busy", "true"));
    expect(button).toHaveFocus();
    await user.keyboard("{Enter}");
    expect(api.add).toHaveBeenCalledTimes(1);
    await act(async () => reject(new Error("Network unavailable")));
    await waitFor(() => expect(button).not.toHaveAttribute("aria-busy"));
    expect(button).toHaveAttribute("aria-pressed", "false");
    expect(button).toHaveAccessibleDescription("3 likes.");
    expect(button).toHaveFocus();
    expect(toast.error).toHaveBeenCalled();
  });
});
