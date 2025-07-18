/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */

/**
 * Parameters for repositioning a node in the hierarchy. You may change the
node's parent using `parent`, and/or reposition it among its siblings
using one of: `before`, `after`, or `index`. Using multiple reordering
properties is not allowed.

 */
export interface NodePositionMutableProps {
  /** Move this node after the sibling with this ID. */
  after?: string;
  /** Move this node before the sibling with this ID. */
  before?: string;
  /**
   * Optional new parent node slug. Set to null to move node to the root.

   * @nullable
   */
  parent?: string | null;
}
