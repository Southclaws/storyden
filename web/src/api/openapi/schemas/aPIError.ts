/**
 * Generated by orval v6.28.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { APIErrorMetadata } from "./aPIErrorMetadata";

/**
 * A description of an error including a human readable message and any
related metadata from the request and associated services.

 */
export interface APIError {
  /** The internal error, not intended for end-user display. */
  error: string;
  /** A human-readable message intended for end-user display. */
  message?: string;
  /** Any additional metadata related to the error. */
  metadata?: APIErrorMetadata;
  /** A suggested action for the user. */
  suggested?: string;
}
