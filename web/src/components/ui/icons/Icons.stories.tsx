import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import type { ComponentType } from "react";

import { css } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

import * as AddIcons from "./Add";
import * as AdminIcons from "./Admin";
import * as ArchiveIcons from "./Archive";
import * as ArrowIcons from "./Arrow";
import * as AuthorIcons from "./Author";
import * as BanIconIcons from "./BanIcon";
import * as BiometricIcons from "./Biometric";
import * as BlocksIcons from "./Blocks";
import * as BookmarkIcons from "./Bookmark";
import * as BreadcrumbIcons from "./Breadcrumb";
import * as BulletIcons from "./Bullet";
import * as CalendarIcons from "./Calendar";
import * as CancelIcons from "./Cancel";
import * as CardIcons from "./Card";
import * as CategoryIcons from "./Category";
import * as CheckIcons from "./Check";
import * as CheckCircleIcons from "./CheckCircle";
import * as ChevronIcons from "./Chevron";
import * as CloseIcons from "./Close";
import * as CollectionIcons from "./Collection";
import * as ColourIcons from "./Colour";
import * as CopyIcons from "./Copy";
import * as CreateIcons from "./Create";
import * as CreateFolderIcons from "./CreateFolder";
import * as DeleteIcons from "./Delete";
import * as DeletedMemberIcons from "./DeletedMember";
import * as DiscussionIcons from "./Discussion";
import * as DraftIcons from "./Draft";
import * as DragHandleIcons from "./DragHandle";
import * as EditIcons from "./Edit";
import * as EmptyIcons from "./Empty";
import * as ExternalLinkIcons from "./ExternalLink";
import * as HideIconIcons from "./HideIcon";
import * as HomeIcons from "./Home";
import * as InboxIcons from "./Inbox";
import * as InfoIcons from "./Info";
import * as IntelligenceIcons from "./Intelligence";
import * as InvitationIcons from "./Invitation";
import * as LayoutIcons from "./Layout";
import * as LayoutGridIcons from "./LayoutGrid";
import * as LayoutListIcons from "./LayoutList";
import * as LayoutTableIcons from "./LayoutTable";
import * as LibraryIcons from "./Library";
import * as LikeIcons from "./Like";
import * as LinkIcons from "./Link";
import * as LogoutIcons from "./Logout";
import * as MediaIcons from "./Media";
import * as MembersIcons from "./Members";
import * as MenuIcons from "./Menu";
import * as MoreIcons from "./More";
import * as NotificationIcons from "./Notification";
import * as ParticipatingIcons from "./Participating";
import * as PinIcons from "./Pin";
import * as ProfileIcons from "./Profile";
import * as QueueIcons from "./Queue";
import * as ReactionIcons from "./Reaction";
import * as ReinstateIconIcons from "./ReinstateIcon";
import * as RemoveIcons from "./Remove";
import * as ReplyIcons from "./Reply";
import * as ReportIcons from "./Report";
import * as RobotIcons from "./Robot";
import * as RolesIcons from "./Roles";
import * as SaveIcons from "./Save";
import * as SearchIcons from "./Search";
import * as SelectIcons from "./Select";
import * as SettingsIcons from "./Settings";
import * as ShareIcons from "./Share";
import * as ShowIconIcons from "./ShowIcon";
import * as SizeIcons from "./Size";
import * as SlugIcons from "./Slug";
import * as SortIcons from "./Sort";
import * as SubmenuIcons from "./Submenu";
import * as SuggestIcons from "./Suggest";
import * as SystemIcons from "./System";
import * as TagIcons from "./Tag";
import * as ToolIcons from "./Tool";
import * as ToolsetIcons from "./Toolset";
import * as TypographyIcons from "./Typography";
import * as VersionsIcons from "./Versions";
import * as VisibilityIcons from "./Visibility";
import * as WarningIcons from "./Warning";

const modules = {
  AddIcons,
  AdminIcons,
  ArchiveIcons,
  ArrowIcons,
  AuthorIcons,
  BanIconIcons,
  BiometricIcons,
  BlocksIcons,
  BookmarkIcons,
  BreadcrumbIcons,
  BulletIcons,
  CalendarIcons,
  CancelIcons,
  CardIcons,
  CategoryIcons,
  CheckIcons,
  CheckCircleIcons,
  ChevronIcons,
  CloseIcons,
  CollectionIcons,
  ColourIcons,
  CopyIcons,
  CreateIcons,
  CreateFolderIcons,
  DeleteIcons,
  DeletedMemberIcons,
  DiscussionIcons,
  DraftIcons,
  DragHandleIcons,
  EditIcons,
  EmptyIcons,
  ExternalLinkIcons,
  HideIconIcons,
  HomeIcons,
  InboxIcons,
  InfoIcons,
  IntelligenceIcons,
  InvitationIcons,
  LayoutIcons,
  LayoutGridIcons,
  LayoutListIcons,
  LayoutTableIcons,
  LibraryIcons,
  LikeIcons,
  LinkIcons,
  LogoutIcons,
  MediaIcons,
  MembersIcons,
  MenuIcons,
  MoreIcons,
  NotificationIcons,
  ParticipatingIcons,
  PinIcons,
  ProfileIcons,
  QueueIcons,
  ReactionIcons,
  ReinstateIconIcons,
  RemoveIcons,
  ReplyIcons,
  ReportIcons,
  RobotIcons,
  RolesIcons,
  SaveIcons,
  SearchIcons,
  SelectIcons,
  SettingsIcons,
  ShareIcons,
  ShowIconIcons,
  SizeIcons,
  SlugIcons,
  SortIcons,
  SubmenuIcons,
  SuggestIcons,
  SystemIcons,
  TagIcons,
  ToolIcons,
  ToolsetIcons,
  TypographyIcons,
  VersionsIcons,
  VisibilityIcons,
  WarningIcons,
} satisfies Record<string, Record<string, unknown>>;

const icons = Object.entries(modules)
  .flatMap(([moduleName, module]) =>
    Object.entries(module)
      .filter(([name, value]) => name.endsWith("Icon") && value)
      .map(([name, Icon]) => ({
        id: `${moduleName}:${name}`,
        name,
        Icon: Icon as ComponentType<any>,
      })),
  )
  .sort((a, b) => a.name.localeCompare(b.name));

const iconPreviewClassName = css({
  "& svg": {
    width: "[100%]",
    height: "[100%]",
    maxWidth: "[100%]",
    maxHeight: "[100%]",
  },
});

const meta = {
  title: "Foundations/Icons/Library",
  parameters: {
    docs: {
      description: {
        component:
          "Storyden icons are semantic product exports, not a shape catalogue. Product code imports from `components/ui/icons`; raw Lucide or Heroicons imports belong inside this library. Choose an icon by its exported meaning—`CollectionIcon` represents Storyden Collections and `ToolsetIcon` represents Robot Toolsets. If no semantic icon matches, prefer no icon before adding decoration. Add a new export only for a recurring product concept. The containing component owns contextual size whenever possible: compact inline and Card title icons render at 16px, while icon buttons derive icon size from the control.",
      },
    },
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const All: Story = {
  render: () => (
    <styled.div
      display="grid"
      gridTemplateColumns="repeat(auto-fill, minmax(9rem, 1fr))"
      gap="3"
    >
      {icons.map(({ id, name, Icon }) => (
        <styled.div
          key={id}
          display="flex"
          alignItems="center"
          gap="2"
          minW="0"
          borderWidth="thin"
          borderColor="border.muted"
          borderRadius="md"
          p="2"
        >
          <styled.span
            color="text.default"
            flexShrink="0"
            display="inline-flex"
            alignItems="center"
            justifyContent="center"
            width="5"
            height="5"
            overflow="hidden"
            className={iconPreviewClassName}
          >
            <Icon width="100%" height="100%" />
          </styled.span>
          <styled.span
            color="text.subtle"
            fontSize="xs"
            overflow="hidden"
            textOverflow="ellipsis"
            whiteSpace="nowrap"
          >
            {name}
          </styled.span>
        </styled.div>
      ))}
    </styled.div>
  ),
};
