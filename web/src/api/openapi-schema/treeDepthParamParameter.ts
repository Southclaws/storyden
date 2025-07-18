/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */

/**
 * When set to a positive value, the nodes in the response will include all
child nodes up to the specified depth. When set to zero, then if the
request includes a node ID only that node will be returned, otherwise
only top-level (root) nodes will be returned.

 */
export type TreeDepthParamParameter = string;
