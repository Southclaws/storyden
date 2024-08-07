/**
 * Generated by orval v6.30.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { Identifier } from "./identifier";
import type { Metadata } from "./metadata";
import type { PostContent } from "./postContent";
import type { Url } from "./url";

export interface PostInitialProps {
  body: PostContent;
  meta?: Metadata;
  reply_to?: Identifier;
  url?: Url;
}
