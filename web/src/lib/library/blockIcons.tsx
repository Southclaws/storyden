import {
  GalleryThumbnailsIcon,
  ImagesIcon,
  LayersIcon,
  TablePropertiesIcon,
  TagsIcon,
  TextIcon,
  TypeIcon,
} from "lucide-react";

import { LinkIcon } from "@/components/ui/icons/Link";

import { LibraryPageBlockType } from "./metadata";

export const LibraryPageBlockIcon: Record<
  LibraryPageBlockType,
  React.ComponentType
> = {
  title: TypeIcon,
  cover: GalleryThumbnailsIcon,
  link: LinkIcon,
  content: TextIcon,
  assets: ImagesIcon,
  properties: TablePropertiesIcon,
  directory: LayersIcon,
  tags: TagsIcon,
};

export function BlockIcon({ blockType }: { blockType: LibraryPageBlockType }) {
  const Icon = LibraryPageBlockIcon[blockType];

  return <Icon />;
}
