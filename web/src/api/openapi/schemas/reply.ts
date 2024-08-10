/**
 * Generated by orval v6.30.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { CommonProperties } from "./commonProperties";
import type { PostProps } from "./postProps";
import type { ReplyProps } from "./replyProps";

/**
 * A new post within a thread of posts. A post may reply to another post in
the thread by specifying the `reply_to` property. The identifier in the
`reply_to` value must be post within the same thread.

 */
export type Reply = CommonProperties & PostProps & ReplyProps;
