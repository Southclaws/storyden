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

export function hasImageFile(items: DataTransferItemList | DataTransferItem[]): boolean {
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
