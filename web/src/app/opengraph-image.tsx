import { streamMetadataImage } from "@/lib/metadata/image-response";

export const size = { width: 1200, height: 630 };
export const contentType = "image/png";
export const dynamic = "force-dynamic";

export default async function Image() {
  return streamMetadataImage("/info/banner");
}
