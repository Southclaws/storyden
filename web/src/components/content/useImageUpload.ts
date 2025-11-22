import { assetUpload } from "src/api/openapi-client/assets";

import { getAssetUploadMutationKey } from "@/api/openapi-client/assets";
import { APIError, Asset, AssetUploadParams } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";
import { deriveError } from "@/utils/error";

export function useImageUpload() {
  async function upload(f: File, params?: AssetUploadParams) {
    if (!isSupportedImage(f.type)) {
      throw new Error(`Unsupported image format ${f.type}`);
    }

    const asset = await assetUpload(f, params);

    return asset;
  }

  async function uploadWithProgress(
    f: File,
    onProgress: (progress: number) => void,
    params?: AssetUploadParams,
    abortController?: AbortController,
  ): Promise<Asset> {
    if (!isSupportedImage(f.type)) {
      throw new Error(`Unsupported image format ${f.type}`);
    }

    const url = buildAssetUploadURL(params);

    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      xhr.withCredentials = true;

      if (abortController) {
        abortController.signal.addEventListener("abort", () => {
          xhr.abort();
        });
      }

      xhr.upload.addEventListener("progress", (event) => {
        if (event.lengthComputable) {
          const progress = (event.loaded / event.total) * 100;
          onProgress(progress);
        }
      });

      xhr.addEventListener("load", () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            // TODO: Validate this properly, remove "as".
            const asset = JSON.parse(xhr.responseText) as Asset;
            resolve(asset);
          } catch (error) {
            reject(new Error("Failed to parse upload response"));
          }
        } else {
          try {
            // TODO: Validate this properly, remove "as".
            const apiError = JSON.parse(xhr.responseText) as APIError;

            // NOTE: Derive error should also handle these cases.
            reject(
              new Error(
                deriveError(
                  apiError.message || apiError.error || "Upload failed",
                ),
              ),
            );
          } catch {
            reject(new Error(`Upload failed with status ${xhr.status}`));
          }
        }
      });

      xhr.addEventListener("error", () => {
        reject(new Error("Network error during upload"));
      });

      xhr.addEventListener("abort", () => {
        reject(new Error("Upload cancelled"));
      });

      xhr.open("POST", url);
      xhr.setRequestHeader("Content-Type", "application/octet-stream");

      // Send the file directly as the body
      xhr.send(f);
    });
  }

  return {
    upload,
    uploadWithProgress,
  };
}

// TODO: Support uploading arbitrary files.
export function isSupportedImage(mime: string): boolean {
  const category = mime.split("/")[0];

  switch (category) {
    case "image":
      return true;

    default:
      return false;
  }
}

export function hasImageFile(
  items: DataTransferItemList | DataTransferItem[],
): boolean {
  const itemArray = Array.from(items);
  return itemArray.some((item) => {
    if ("kind" in item && item.kind !== "file") {
      return false;
    }
    return item.type.startsWith("image/");
  });
}

export function getImageFiles(files: FileList | File[]): File[] {
  return Array.from(files).filter((file) => isSupportedImage(file.type));
}

function buildAssetUploadURL(params?: AssetUploadParams): URL {
  const [uploadPath, uploadParams] = getAssetUploadMutationKey(params);

  const query = uploadParams
    ? "?" + new URLSearchParams(uploadParams).toString()
    : "";

  const path = "/api" + uploadPath + query;

  const url = new URL(path, API_ADDRESS);

  return url;
}
