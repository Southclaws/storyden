/**
 * Generated by orval v6.9.6 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { PostsCreateOKResponse, PostsCreateBody } from "./schemas";
import { fetcher } from "../client";

/**
 * Create a new post within a thread.
 */
export const postsCreate = (
  threadMark: string,
  postsCreateBody: PostsCreateBody
) => {
  return fetcher<PostsCreateOKResponse>({
    url: `/v1/threads/${threadMark}/posts`,
    method: "post",
    headers: { "Content-Type": "application/json" },
    data: postsCreateBody,
  });
};
