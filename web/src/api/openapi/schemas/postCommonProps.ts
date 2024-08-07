/**
 * Generated by orval v6.30.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AssetList } from "./assetList";
import type { Identifier } from "./identifier";
import type { LinkList } from "./linkList";
import type { Metadata } from "./metadata";
import type { PostContent } from "./postContent";
import type { ProfileReference } from "./profileReference";
import type { ReactList } from "./reactList";
import type { ThreadMark } from "./threadMark";

export interface PostCommonProps {
  assets: AssetList;
  author: ProfileReference;
  body: PostContent;
  links: LinkList;
  meta?: Metadata;
  reacts: ReactList;
  reply_to?: Identifier;
  root_id: Identifier;
  root_slug: ThreadMark;
}
