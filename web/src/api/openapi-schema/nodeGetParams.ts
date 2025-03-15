/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type { NodeChildrenSortParamParameter } from "./nodeChildrenSortParamParameter";
import type { PaginationQueryParameter } from "./paginationQueryParameter";

export type NodeGetParams = {
  /**
 * The field (either in schema or in property schema) to sort by.

 */
  children_sort?: NodeChildrenSortParamParameter;
  /**
   * Pagination query parameters.
   */
  page?: PaginationQueryParameter;
};
