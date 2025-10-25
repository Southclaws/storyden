// Curated MIME type → extensions map for common upload formats
// Covers image, video, and document types we accept
export const MIME_EXTENSIONS: Record<string, string[]> = {
  // Images
  "image/jpeg": ["jpg", "jpeg"],
  "image/png": ["png"],
  "image/gif": ["gif"],
  "image/svg+xml": ["svg"],
  "image/webp": ["webp"],
  "image/avif": ["avif"],
  "image/bmp": ["bmp"],
  "image/tiff": ["tiff", "tif"],

  // Videos
  "video/mp4": ["mp4"],
  "video/webm": ["webm"],
  "video/quicktime": ["mov"],
  "video/x-msvideo": ["avi"],
  "video/3gpp": ["3gp"],
  "video/x-flv": ["flv"],

  // Documents
  "application/pdf": ["pdf"],
  "application/msword": ["doc"],
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document": [
    "docx",
  ],
  "application/vnd.ms-excel": ["xls"],
  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ["xlsx"],
  "application/vnd.ms-powerpoint": ["ppt"],
  "application/vnd.openxmlformats-officedocument.presentationml.presentation": [
    "pptx",
  ],
  "text/plain": ["txt"],
  "text/csv": ["csv"],
  "application/rtf": ["rtf"],

  // Archives
  "application/zip": ["zip"],
  "application/x-rar-compressed": ["rar"],
  "application/x-7z-compressed": ["7z"],
  "application/gzip": ["gz"],
  "application/x-tar": ["tar"],
};

/**
 * Get file extensions for a given MIME type
 */
export function getExtensionsForMimeType(mimeType: string): string[] {
  return MIME_EXTENSIONS[mimeType] || [];
}

/**
 * Get all extensions for a list of MIME types
 */
export function getExtensionsForMimeTypes(mimeTypes: string[]): string[] {
  return mimeTypes.reduce((prev: string[], mimeType: string) => {
    const extensions = getExtensionsForMimeType(mimeType);
    return [...prev, ...extensions];
  }, []);
}
