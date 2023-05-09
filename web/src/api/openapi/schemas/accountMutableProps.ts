/**
 * Generated by orval v6.14.3 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AccountHandle } from "./accountHandle";
import type { AccountName } from "./accountName";
import type { TagList } from "./tagList";

export interface AccountMutableProps {
  handle?: AccountHandle;
  name?: AccountName;
  bio?: string;
  interests?: TagList;
}
