import { streamMetadataImage } from "@/lib/metadata/image-response";

export const size = {
  width: 32,
  height: 32,
};
export const contentType = "image/png";
export const dynamic = "force-dynamic";

export default async function Icon() {
  return streamMetadataImage("/info/icon/32x32");
}
