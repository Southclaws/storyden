import { fireEvent, render, screen } from "@testing-library/react";
import type { ComponentProps } from "react";
import { describe, expect, it, vi } from "vitest";

import { useReplyGet } from "@/api/openapi-client/replies";
import type { Reply, Thread } from "@/api/openapi-schema";
import { RenderedContentProvider } from "@/components/content/ContentRenderer/RenderedContentContext";

import { ThreadConclusion } from "./ThreadConclusion";

const refresh = vi.fn();

function queryResult(
  data?: Reply,
  error?: unknown,
): ReturnType<typeof useReplyGet> {
  return {
    data,
    error,
    mutate: refresh,
    isValidating: false,
    isLoading: false,
    swrKey: ["conclusion-test"],
  };
}
vi.mock("@/api/openapi-client/replies", () => ({ useReplyGet: vi.fn() }));
vi.mock("@/auth", () => ({ useSession: () => undefined }));
vi.mock("next/navigation", () => ({
  useRouter: () => ({ prefetch: vi.fn() }),
}));

const answer = {
  id: "answer",
  root_id: "thread-id",
  body: "<p>The selected answer</p>",
  visibility: "published",
  title: "A discussion",
  slug: "thread-slug#answer",
  root_slug: "thread-slug",
  assets: [],
  reacts: [],
  likes: { likes: 0, liked: false },
  collections: { has_collected: false, in_collections: 0 },
  body_links: [],
  author: {
    id: "member",
    handle: "member",
    name: "Member",
    roles: [],
    joined: "2026-11-08T10:00:00Z",
  },
  createdAt: "2026-11-08T10:00:00Z",
  updatedAt: "2026-11-08T10:00:00Z",
} as Reply;
const thread = {
  id: "thread-id",
  slug: "thread-slug",
  conclusion_reply_id: answer.id,
  replies: { replies: [] },
} as unknown as Thread;

function Preview(props: ComponentProps<typeof ThreadConclusion>) {
  return (
    <RenderedContentProvider
      content={{ [answer.body]: <p>The selected answer</p> }}
    >
      <ThreadConclusion {...props} />
    </RenderedContentProvider>
  );
}

describe("thread conclusion preview", () => {
  it("renders the SSR answer and links to its original position on another page", () => {
    const fetch = vi.spyOn(globalThis, "fetch");

    vi.mocked(useReplyGet).mockImplementation((_slug, _id, options) =>
      queryResult(options?.swr?.fallbackData as Reply | undefined),
    );
    render(<Preview thread={thread} initialConclusionReply={answer} />);
    expect(screen.getByRole("article", { name: "Answer" })).toBeInTheDocument();
    expect(screen.getByText("The selected answer")).toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: "View in thread" }),
    ).toHaveAttribute("href", "/t/locate/answer");
    expect(fetch).not.toHaveBeenCalled();
  });

  it("renders a pointer instead of duplicating a reply on the current page", () => {
    vi.mocked(useReplyGet).mockReturnValue(queryResult(answer));
    render(
      <Preview
        thread={{
          ...thread,
          replies: { ...thread.replies, replies: [answer] },
        }}
        initialConclusionReply={answer}
      />,
    );
    expect(screen.queryByText("The selected answer")).not.toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: "Jump to answer" }),
    ).toHaveAttribute("href", "#answer");
    expect(vi.mocked(useReplyGet).mock.lastCall?.[2]?.swr?.enabled).toBe(false);
  });

  it("never renders a previous answer while loading a replacement", () => {
    vi.mocked(useReplyGet).mockReturnValue(queryResult(answer));
    const { rerender } = render(
      <Preview thread={thread} initialConclusionReply={answer} />,
    );
    rerender(
      <Preview
        thread={{ ...thread, conclusion_reply_id: "replacement" }}
        initialConclusionReply={answer}
      />,
    );
    expect(screen.queryByText("The selected answer")).not.toBeInTheDocument();
    expect(screen.getByRole("status")).toHaveTextContent("Loading answer");
    expect(
      vi.mocked(useReplyGet).mock.lastCall?.[2]?.swr?.fallbackData,
    ).toBeUndefined();
  });

  it("hides the preview when the conclusion is removed or the thread is reopened", () => {
    vi.mocked(useReplyGet).mockReturnValue(queryResult(answer));
    render(
      <Preview
        thread={{ ...thread, conclusion_reply_id: undefined }}
        initialConclusionReply={answer}
      />,
    );
    expect(screen.queryByRole("article")).not.toBeInTheDocument();
    expect(vi.mocked(useReplyGet).mock.lastCall?.[2]?.swr?.enabled).toBe(false);
  });

  it("offers recovery when loading fails", () => {
    vi.mocked(useReplyGet).mockReturnValue(
      queryResult(undefined, new Error("Unavailable")),
    );
    render(<Preview thread={thread} />);
    expect(screen.getByRole("status")).toHaveTextContent(
      "The answer could not be loaded",
    );
    fireEvent.click(screen.getByRole("button", { name: "Try again" }));
    expect(refresh).toHaveBeenCalledOnce();
  });
});
