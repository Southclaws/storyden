import { act, render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { toast } from "sonner";
import useSWR, { SWRConfig } from "swr";
import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  getThreadGetKey,
  getThreadListKey,
} from "@/api/openapi-client/threads";
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
  slug: "harbour-notes",
  likes: { likes: 3, liked: false },
} as ThreadReference;
let serverPage: ThreadListOKResponse;

function Feed({
  signedIn = true,
  detail = false,
}: {
  signedIn?: boolean;
  detail?: boolean;
}) {
  const [initialPage] = useState(serverPage);
  const { data } = useSWR(
    getThreadListKey({ page: "1" }),
    async () => serverPage,
    { fallbackData: initialPage },
  );
  const { data: detailThread } = useSWR(
    detail ? getThreadGetKey(thread.slug, { page: "2" }) : null,
    async () => serverPage.threads[0]!,
    { fallbackData: initialPage.threads[0]! },
  );

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        throw new Error("Liking must not submit the thread editor");
      }}
    >
      <LikeButton
        thread={detail ? detailThread! : data!.threads[0]!}
        initialSession={signedIn ? account : undefined}
        label="Like Harbour notes"
        size={detail ? "large" : "small"}
      />
    </form>
  );
}

function renderFeed(signedIn = true, detail = false) {
  return render(<Feed signedIn={signedIn} detail={detail} />, {
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

describe.each([false, true])("likes on thread detail: %s", (detail) => {
  beforeEach(() => {
    serverPage = { threads: [thread] } as ThreadListOKResponse;
  });

  it("reports likes and the sign-in requirement for guests", () => {
    renderFeed(false, detail);

    const button = screen.getByRole("button", { name: "Like Harbour notes" });
    expect(button).toBeDisabled();
    expect(button).toHaveAccessibleDescription(
      "3 likes. Sign in to like discussions.",
    );

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
    renderFeed(true, detail);
    const button = screen.getByRole("button", {
      name: "Like Harbour notes",
      pressed: false,
    });
    expect(button).toHaveAccessibleDescription("3 likes.");

    if (detail) {
      expect(within(button).getByText("Like")).toBeVisible();
    } else {
      expect(within(button).queryByText("Like")).not.toBeInTheDocument();
    }

    expect(within(button).getByText("3")).toBeVisible();

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
    renderFeed(true, detail);
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
