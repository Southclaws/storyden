/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: rolling
 */
import type { AssetIDs } from "./assetIDs";
import type { AssetSourceList } from "./assetSourceList";
import type { Identifier } from "./identifier";
import type { Metadata } from "./metadata";
import type { NodeName } from "./nodeName";
import type { PostContent } from "./postContent";
import type { PropertyList } from "./propertyList";
import type { Slug } from "./slug";
import type { TagNameList } from "./tagNameList";
import type { Url } from "./url";
import type { Visibility } from "./visibility";

export interface NodeInitialProps {
  asset_ids?: AssetIDs;
  asset_sources?: AssetSourceList;
  content?: PostContent;
  meta?: Metadata;
  name: NodeName;
  parent?: Slug;
  primary_image_asset_id?: Identifier;
  properties?: PropertyList;
  slug?: Slug;
  tags?: TagNameList;
  url?: Url;
  visibility?: Visibility;
}
