"use client";

import { Presence } from "@ark-ui/react";

import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";

import { HStack } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";

import { CloseAction } from "../../Action/Close";
import { AdminAction } from "../Anchors/Admin";
import { ComposeAction } from "../Anchors/Compose";
import { DashboardAction } from "../Anchors/Dashboard";
import { DraftsAction } from "../Anchors/Drafts";
import { HomeAction } from "../Anchors/Home";
import { LibraryAction } from "../Anchors/Library";
import { LoginAction } from "../Anchors/Login";
import { LogoutAction } from "../Anchors/Logout";
import { SettingsAction } from "../Anchors/Settings";
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
                <HomeAction />
                <DraftsAction />
                <LogoutAction />
              </HStack>
              <HStack>
                {account.admin && (
                  <>
                    <AdminAction />
                    {/* TODO: Move public drafts for admin review to /queue */}
                    {/* <QueueAction /> */}
                  </>
                )}
                <SettingsAction />
              </HStack>
            </>
          ) : (
            <LoginAction />
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
              <HomeAction />
              <ComposeAction />
              <LibraryAction />
              <DashboardAction onClick={onExpand} />
            </>
          )}
        </HStack>
      ) : (
        <HStack gap="4" w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <HomeAction />
              <Search />
              <CloseAction onClick={onClose} />
            </>
          ) : (
            <>
              <HomeAction />
              <LoginAction />
              <DashboardAction onClick={onExpand} />
            </>
          )}
        </HStack>
      )}
    </Toolpill>
  );
}
