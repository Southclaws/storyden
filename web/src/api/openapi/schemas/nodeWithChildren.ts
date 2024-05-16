/**
 * Generated by orval v6.28.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { DatagraphRecommendations } from "./datagraphRecommendations";
import type { Node } from "./node";
import type { NodeWithChildrenAllOf } from "./nodeWithChildrenAllOf";

/**
 * The full properties of a node including all child nodes.

 */
export type NodeWithChildren = Node &
  DatagraphRecommendations &
  NodeWithChildrenAllOf;
