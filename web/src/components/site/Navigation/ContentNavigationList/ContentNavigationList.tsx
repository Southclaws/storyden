"use client";

import {
  type Account,
  type CategoryListOKResponse,
  InstanceCapability,
  type NodeListResult,
} from "@/api/openapi-schema";
import { hasCapability } from "@/lib/settings/capabilities";
import { useNavigationConfig } from "@/lib/settings/navigation-client";
import { type Settings } from "@/lib/settings/settings";
import { useSettings } from "@/lib/settings/settings-client";
import { useSiteEditorState } from "@/lib/settings/site-editor-client";
import { styled } from "@/styled-system/jsx";
import { useOverflowGradient } from "@/utils/useOverflowGradient";

import { useNavigation } from "../useNavigation";

import { NavigationItems } from "./NavigationItems";

type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
};

export function ContentNavigationList(props: Props) {
  const { currentPath } = useNavigation();
  const scrollViewportRef = useOverflowGradient<HTMLDivElement>();
  const navigation = useNavigationConfig(props.initialSettings, false);
  const settings = useSettings(props.initialSettings, false);
  const robotsEnabled =
    settings.ready &&
    hasCapability(InstanceCapability.robots, settings.settings.capabilities);
  const { isEditing } = useSiteEditorState({
    initialSession: props.initialSession,
    initialSettings: props.initialSettings,
  });

  return (
    <styled.nav
      aria-label="Site navigation"
      display="flex"
      flexDir="column"
      height="full"
      width="full"
      minH="0"
      alignItems="start"
    >
      <div
        ref={scrollViewportRef}
        className="navigation__scroll-viewport navigation-editor__viewport"
        style={{ scrollbarWidth: "none" }}
      >
        <NavigationItems
          navigation={navigation}
          currentPath={currentPath}
          initialCategoryList={props.initialCategoryList}
          initialNodeList={props.initialNodeList}
          isEditing={isEditing}
          robotsEnabled={robotsEnabled}
        />
      </div>
    </styled.nav>
  );
}
