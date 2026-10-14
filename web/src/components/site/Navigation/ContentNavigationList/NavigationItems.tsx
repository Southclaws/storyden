import { Suspense, lazy } from "react";

import {
  type Account,
  type CategoryListOKResponse,
  type NodeListResult,
} from "@/api/openapi-schema";
import {
  type NavigationConfig,
  getNavigationItemKey,
} from "@/lib/settings/navigation";

import { NavigationItemContent } from "./NavigationItemContent";

export type Props = {
  navigation: NavigationConfig;
  currentPath: string;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
  session?: Account;
  isEditing: boolean;
  robotsEnabled: boolean;
};

const NavigationItemsEditor = lazy(() =>
  import("./NavigationItemsEditor").then((module) => ({
    default: module.NavigationItemsEditor,
  })),
);

export function NavigationItems(props: Props) {
  if (!props.isEditing) return <NavigationItemsView {...props} />;
  return (
    <Suspense fallback={<NavigationItemsView {...props} />}>
      <NavigationItemsEditor {...props} />
    </Suspense>
  );
}

function NavigationItemsView(props: Props) {
  const visibleItems = props.navigation.items.filter(
    (item) => item.type !== "robots" || props.robotsEnabled,
  );
  return (
    <div className="navigation-editor__items" data-editing={false}>
      {visibleItems.map((item) => (
        <div
          key={getNavigationItemKey(item)}
          className="navigation-editor__item-content"
        >
          <NavigationItemContent
            item={item}
            currentPath={props.currentPath}
            session={props.session}
            initialCategoryList={props.initialCategoryList}
            initialNodeList={props.initialNodeList}
          />
        </div>
      ))}
    </div>
  );
}
