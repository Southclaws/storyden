import { describe, expect, it, vi } from "vitest";

import { RequestError } from "@/api/common";
import type { Reply, Thread } from "@/api/openapi-schema";
import { replyGet } from "@/api/openapi-server/replies";

import { getConclusionReply } from "./conclusion-server";

vi.mock("server-only", () => ({}));
vi.mock("@/api/openapi-server/replies", () => ({ replyGet: vi.fn() }));

const answer = {
  id: "answer",
  root_id: "thread-id",
  body: "<p>The answer</p>",
} as Reply;
const thread = {
  id: "thread-id",
  slug: "thread-slug",
  conclusion_reply_id: answer.id,
  replies: { replies: [] },
} as unknown as Thread;

describe("server conclusion reply", () => {
  it("does not fetch when no conclusion is selected", async () => {
    expect(
      await getConclusionReply({ ...thread, conclusion_reply_id: undefined }),
    ).toBeUndefined();
    expect(replyGet).not.toHaveBeenCalled();
  });

  it("reuses an answer already fetched in the server's reply page", async () => {
    expect(
      await getConclusionReply({
        ...thread,
        replies: { ...thread.replies, replies: [answer] },
      }),
    ).toEqual(answer);
    expect(replyGet).not.toHaveBeenCalled();
  });

  it("fetches the answer independently of the current reply page", async () => {
    vi.mocked(replyGet).mockResolvedValue({
      data: answer,
      status: 200,
      headers: new Headers(),
    });
    expect(await getConclusionReply(thread)).toEqual(answer);
    expect(replyGet).toHaveBeenCalledWith(thread.slug, answer.id);
  });

  it.each([403, 404])(
    "handles a conclusion becoming unavailable during SSR (%i)",
    async (status) => {
      vi.mocked(replyGet).mockRejectedValue(
        new RequestError("Unavailable", status),
      );
      expect(await getConclusionReply(thread)).toBeUndefined();
    },
  );

  it("preserves operational errors", async () => {
    vi.mocked(replyGet).mockRejectedValue(new Error("API unavailable"));
    await expect(getConclusionReply(thread)).rejects.toThrow("API unavailable");
  });
});
