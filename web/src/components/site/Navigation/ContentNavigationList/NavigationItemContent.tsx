import {
  type Account,
  type CategoryListOKResponse,
  type NodeListResult,
  Permission,
} from "@/api/openapi-schema";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { ExternalLinkIcon } from "@/components/ui/icons/ExternalLink";
import {
  BuiltInNavigationItemName,
  type NavigationItem,
} from "@/lib/settings/navigation";
import { hasPermission } from "@/utils/permissions";

import { Anchor } from "../Anchors/Anchor";
import { CollectionsAnchor } from "../Anchors/Collections";
import { ComposeAnchor } from "../Anchors/Compose";
import { LinksAnchor } from "../Anchors/Link";
import { MembersAnchor } from "../Anchors/Members";
import { RobotsAnchor } from "../Anchors/Robots";
import { RolesAnchor } from "../Anchors/Roles";
import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";
import { libraryNavigationVisibility } from "../LibraryNavigationTree/LibraryNavigationTree.constants";

export function NavigationItemContent({
  item,
  currentPath,
  session,
  initialNodeList,
  initialCategoryList,
}: {
  item: NavigationItem;
  currentPath: string;
  session?: Account;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
}) {
  switch (item.type) {
    case "compose":
      return <ComposeAnchor label={BuiltInNavigationItemName.compose} />;
    case "categories":
      return (
        <CategoryList
          initialCategoryList={initialCategoryList}
          initialSession={session}
          currentPath={currentPath}
        />
      );
    case "library":
      return (
        <LibraryNavigationTree
          initialNodeList={initialNodeList}
          initialSession={session}
          currentPath={currentPath}
          visibility={libraryNavigationVisibility}
        />
      );
    case "robots":
      return (
        <RobotsAnchor
          route={
            hasPermission(
              session,
              Permission.USE_ROBOTS,
              Permission.MANAGE_ROBOTS,
            )
              ? "/robots"
              : "/robots/trails"
          }
        />
      );
    case "collections":
      return <CollectionsAnchor />;
    case "links":
      return <LinksAnchor />;
    case "members":
      return <MembersAnchor />;
    case "roles":
      return <RolesAnchor />;
    case "custom-link":
      return (
        <Anchor
          id={`custom-link-${item.id}`}
          icon={<ExternalLinkIcon />}
          label={item.label}
          route={item.href}
        />
      );
  }
}
