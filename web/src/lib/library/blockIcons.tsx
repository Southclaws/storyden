import {
  AssetsBlockIcon,
  ContentBlockIcon,
  CoverBlockIcon,
  DirectoryBlockIcon,
  LinkBlockIcon,
  PropertiesBlockIcon,
  TagsBlockIcon,
  TitleBlockIcon,
} from "@/components/ui/icons/Blocks";

import { LibraryPageBlockType } from "./metadata";

export const LibraryPageBlockIcon: Record<
  LibraryPageBlockType,
  React.ComponentType
> = {
  title: TitleBlockIcon,
  cover: CoverBlockIcon,
  link: LinkBlockIcon,
  content: ContentBlockIcon,
  assets: AssetsBlockIcon,
  properties: PropertiesBlockIcon,
  directory: DirectoryBlockIcon,
  tags: TagsBlockIcon,
};

export function BlockIcon({ blockType }: { blockType: LibraryPageBlockType }) {
  const Icon = LibraryPageBlockIcon[blockType];

  return <Icon />;
}
