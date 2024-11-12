/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type { AccountHandleQueryParamParameter } from "./accountHandleQueryParamParameter";
import type { CollectionHasItemQueryParamParameter } from "./collectionHasItemQueryParamParameter";

export type CollectionListParams = {
  /**
   * Account handle.
   */
  account_handle?: AccountHandleQueryParamParameter;
  /**
 * When specified, will include a field in the response indicating whether
or not the specified item is present in the collection. This saves you
needing to make two queries to check if an item is in a collection.

 */
  has_item?: CollectionHasItemQueryParamParameter;
};
