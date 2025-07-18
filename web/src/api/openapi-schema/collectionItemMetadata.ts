/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import type { CollectionItemMembershipType } from "./collectionItemMembershipType";
import type { ProfileReference } from "./profileReference";
import type { RelevanceScore } from "./relevanceScore";

export interface CollectionItemMetadata {
  /** The time that the item was added to the collection. */
  added_at: string;
  membership_type: CollectionItemMembershipType;
  owner: ProfileReference;
  relevance_score?: RelevanceScore;
}
