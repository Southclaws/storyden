import {
  type Account,
  type CategoryListOKResponse,
  type NodeListResult,
} from "@/api/openapi-schema";
import { categoryListCached } from "@/lib/category/server-category-list";
import { nodeListCached } from "@/lib/library/server-node-list";
import { type Settings } from "@/lib/settings/settings";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

type NavigationPaneProps = {
  initialSession?: Account;
  initialSettings?: Settings;
};

export async function NavigationPane({
  initialSession,
  initialSettings,
}: NavigationPaneProps) {
  try {
    const [{ data: initialNodeList }, { data: initialCategoryList }] =
      await Promise.all([
        nodeListCached({
          // NOTE: This doesn't work due to a bug in Orval.
          // visibility: ["draft", "review", "unlisted", "published"],
        }),
        categoryListCached(),
      ]);

    return (
      <NavigationPaneContent
        initialNodeList={initialNodeList}
        initialCategoryList={initialCategoryList}
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
    );
  } catch {
    return (
      <NavigationPaneContent
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
    );
  }
}

type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
};

function NavigationPaneContent({
  initialSession,
  initialSettings,
  initialNodeList,
  initialCategoryList,
}: Props) {
  return (
    <header className="navigation-pane navigation__surface">
      <div id="desktop-nav-box" className="navigation-pane__content">
        <ContentNavigationList
          initialSession={initialSession}
          initialSettings={initialSettings}
          initialNodeList={initialNodeList}
          initialCategoryList={initialCategoryList}
        />
      </div>
    </header>
  );
}
