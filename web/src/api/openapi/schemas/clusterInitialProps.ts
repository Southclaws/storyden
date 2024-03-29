/**
 * Generated by orval v6.24.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AssetIDs } from "./assetIDs";
import type { ClusterDescription } from "./clusterDescription";
import type { ClusterName } from "./clusterName";
import type { PostContent } from "./postContent";
import type { Properties } from "./properties";
import type { Slug } from "./slug";
import type { Url } from "./url";
import type { Visibility } from "./visibility";

export interface ClusterInitialProps {
  asset_ids?: AssetIDs;
  content?: PostContent;
  description: ClusterDescription;
  name: ClusterName;
  parent?: Slug;
  properties?: Properties;
  slug: Slug;
  url?: Url;
  visibility?: Visibility;
}
