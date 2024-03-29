/**
 * Generated by orval v6.24.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AssetList } from "./assetList";
import type { ItemDescription } from "./itemDescription";
import type { ItemName } from "./itemName";
import type { Link } from "./link";
import type { PostContent } from "./postContent";
import type { ProfileReference } from "./profileReference";
import type { Properties } from "./properties";
import type { Slug } from "./slug";
import type { Visibility } from "./visibility";

/**
 * The main properties for an item.
 */
export interface ItemCommonProps {
  assets: AssetList;
  content?: PostContent;
  description: ItemDescription;
  link?: Link;
  name: ItemName;
  owner: ProfileReference;
  properties: Properties;
  slug: Slug;
  visibility: Visibility;
}
