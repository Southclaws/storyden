/**
 * Generated by orval v6.30.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { Visibility } from "./visibility";

/**
 * Filter nodes with specific visibility values. Note that by
default, only published nodes are returned. When 'draft' is
specified, only drafts owned by the requesting account are included.
When 'review' is specified, the request will fail if the requesting
account is not an administrator.

 */
export type VisibilityParamParameter = Visibility[];
