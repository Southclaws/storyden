/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import type { AssetNameQueryParameter } from "./assetNameQueryParameter";
import type { ParentAssetIDQueryParameter } from "./parentAssetIDQueryParameter";

export type AssetUploadParams = {
  /**
   * The client-provided file name for the asset.
   */
  filename?: AssetNameQueryParameter;
  /**
 * For uploading new versions of an existing asset, set this parameter to
the asset ID of the parent asset. This must be an ID and not a filename.
This feature is used for situations where you want to replace an asset
in its usage context, but retain the original with a way to reference it
for features such as editable/croppable images or file version history.

 */
  parent_asset_id?: ParentAssetIDQueryParameter;
};
