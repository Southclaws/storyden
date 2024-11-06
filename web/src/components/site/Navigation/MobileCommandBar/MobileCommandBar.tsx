"use client";

import { Presence } from "@ark-ui/react";

import { Toolpill } from "src/components/site/Toolpill/Toolpill";

import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { HStack, WStack } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { CloseAction } from "../../Action/Close";
import { AccountMenu } from "../AccountMenu/AccountMenu";
import { AdminAnchor } from "../Anchors/Admin";
import { ComposeAnchor } from "../Anchors/Compose";
import { DraftsAnchor } from "../Anchors/Drafts";
import { HomeAnchor } from "../Anchors/Home";
import { LibraryAnchor } from "../Anchors/Library";
import { LoginAnchor } from "../Anchors/Login";
import { LogoutAnchor } from "../Anchors/Logout";
import { SettingsAnchor } from "../Anchors/Settings";
import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";
import { Search } from "../Search/Search";

import { useMobileCommandBar } from "./useMobileCommandBar";

export function MobileCommandBar() {
  const { isExpanded, onExpand, onClose, account } = useMobileCommandBar();
  return (
    <Toolpill onClickOutside={onClose}>
      <Presence
        present={isExpanded}
        className={vstack({ w: "full", gap: "2" })}
      >
        <WStack>
          {account ? (
            <>
              <HStack>
                <HomeAnchor hideLabel />
                <DraftsAnchor hideLabel />
                <LogoutAnchor hideLabel size="xs" />
              </HStack>
              <HStack>
                {account.admin && (
                  <>
                    <AdminAnchor hideLabel />
                    {/* TODO: Move public drafts for admin review to /queue */}
                    {/* <QueueAction /> */}
                  </>
                )}
                <SettingsAnchor hideLabel />
              </HStack>
            </>
          ) : (
            <LoginAnchor />
          )}
        </WStack>

        <ContentNavigationList />
      </Presence>

      {account ? (
        <WStack alignItems="center">
          {isExpanded ? (
            <>
              <AccountMenu account={account} size="sm" />
              <Search />
              <CloseAction onClick={onClose} size="sm" />
            </>
          ) : (
            <>
              <AccountMenu account={account} size="sm" />
              <HomeAnchor hideLabel size="sm" />
              <ComposeAnchor hideLabel size="sm" />
              <LibraryAnchor hideLabel size="sm" />
              <ExpandTrigger onClick={onExpand} />
            </>
          )}
        </WStack>
      ) : (
        <WStack alignItems="center">
          {isExpanded ? (
            <>
              <HomeAnchor hideLabel size="sm" />
              <Search />
              <CloseAction onClick={onClose} size="sm" />
            </>
          ) : (
            <>
              <HomeAnchor hideLabel />
              <LoginAnchor />
              <ExpandTrigger onClick={onExpand} />
            </>
          )}
        </WStack>
      )}
    </Toolpill>
  );
}

function ExpandTrigger(props: ButtonProps) {
  return (
    <IconButton
      title="Main navigation menu"
      variant="ghost"
      size="sm"
      {...props}
    >
      <MenuIcon />
    </IconButton>
  );
}
