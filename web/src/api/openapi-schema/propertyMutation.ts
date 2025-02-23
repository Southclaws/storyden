/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type { Identifier } from "./identifier";
import type { PropertyName } from "./propertyName";
import type { PropertySortKey } from "./propertySortKey";
import type { PropertyType } from "./propertyType";
import type { PropertyValue } from "./propertyValue";

/**
 * A property mutation is a change to a property on a node. It can be used
to update existing properties or add new properties to a node. When a
property already exists by name/fid, type and sort columns are optional.

 */
export interface PropertyMutation {
  fid?: Identifier;
  name: PropertyName;
  sort?: PropertySortKey;
  type?: PropertyType;
  value: PropertyValue;
}
