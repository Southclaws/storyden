"use client";

import { Presence } from "@ark-ui/react";

import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";

import { HStack } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { CloseAction } from "../../Action/Close";
import { MenuAction } from "../Actions/Menu";
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

import { useNavpill } from "./useNavpill";

export function Navpill() {
  const { isExpanded, onExpand, onClose, account } = useNavpill();
  return (
    <Toolpill onClickOutside={onClose}>
      <Presence
        present={isExpanded}
        className={vstack({ w: "full", gap: "2" })}
      >
        <HStack w="full" justify="space-between">
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
        </HStack>

        <ContentNavigationList />
      </Presence>

      {account ? (
        <HStack gap="4" w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <ProfilePill profileReference={account} showHandle={false} />

              <Search />
              <CloseAction onClick={onClose} />
            </>
          ) : (
            <>
              <ProfilePill profileReference={account} showHandle={false} />
              <HomeAnchor hideLabel />
              <ComposeAnchor hideLabel />
              <LibraryAnchor hideLabel />
              <MenuAction onClick={onExpand} />
            </>
          )}
        </HStack>
      ) : (
        <HStack gap="4" w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <HomeAnchor hideLabel />
              <Search />
              <CloseAction onClick={onClose} />
            </>
          ) : (
            <>
              <HomeAnchor hideLabel />
              <LoginAnchor />
              <MenuAction onClick={onExpand} />
            </>
          )}
        </HStack>
      )}
    </Toolpill>
  );
}
