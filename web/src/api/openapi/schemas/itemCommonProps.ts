/**
 * Generated by orval v6.17.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AssetURL } from "./assetURL";
import type { ItemDescription } from "./itemDescription";
import type { ItemName } from "./itemName";
import type { Link } from "./link";
import type { PostContent } from "./postContent";
import type { ProfileReference } from "./profileReference";
import type { Properties } from "./properties";
import type { Slug } from "./slug";

/**
 * The main properties for an item.
 */
export interface ItemCommonProps {
  name: ItemName;
  slug: Slug;
  image_url?: AssetURL;
  link?: Link;
  description: ItemDescription;
  content?: PostContent;
  owner: ProfileReference;
  properties: Properties;
}
