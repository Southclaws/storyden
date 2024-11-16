/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type { AuthMode } from "./authMode";
import type { Metadata } from "./metadata";
import type { PostContent } from "./postContent";

export interface AdminSettingsMutableProps {
  accent_colour?: string;
  authentication_mode?: AuthMode;
  content?: PostContent;
  description?: string;
  /** The settings metadata may be used by frontends to store arbitrary
vendor-specific configuration data specific to the frontend itself.
 */
  metadata?: Metadata;
  title?: string;
}
