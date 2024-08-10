import { assetUpload } from "src/api/openapi-client/assets";

export function useImageUpload() {
  async function upload(f: File) {
    if (!isSupportedImage(f.type)) {
      throw new Error(`Unsupported image format ${f.type}`);
    }

    const asset = await assetUpload(f);

    return asset;
  }

  return {
    upload,
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
