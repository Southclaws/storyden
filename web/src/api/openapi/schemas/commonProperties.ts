/**
 * Generated by orval v6.15.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { Identifier } from "./identifier";
import type { CommonPropertiesMisc } from "./commonPropertiesMisc";

export interface CommonProperties {
  id: Identifier;
  /** The time the resource was created. */
  createdAt: string;
  /** The time the resource was updated. */
  updatedAt: string;
  /** The time the resource was soft-deleted. */
  deletedAt?: string;
  /** Arbitrary extra data stored with the resource. */
  misc?: CommonPropertiesMisc;
}
