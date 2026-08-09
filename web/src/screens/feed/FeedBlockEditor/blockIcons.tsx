import {
  ContentBlockIcon,
  CoverBlockIcon,
  DirectoryBlockIcon,
  TitleBlockIcon,
} from "@/components/ui/icons/Blocks";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { ShareIcon } from "@/components/ui/icons/Share";
import { TextIcon } from "@/components/ui/icons/Typography";
import { FeedBlockType } from "@/lib/settings/feed";

const FeedBlockIcons: Record<FeedBlockType, React.ComponentType> = {
  title: TitleBlockIcon,
  subtitle: TextIcon,
  content: ContentBlockIcon,
  cover: CoverBlockIcon,
  categories: DirectoryBlockIcon,
  threads: DiscussionIcon,
  "quick-share": ShareIcon,
  library: LibraryIcon,
};

export function FeedBlockIcon({ type }: { type: FeedBlockType }) {
  const Icon = FeedBlockIcons[type];
  return <Icon />;
}
