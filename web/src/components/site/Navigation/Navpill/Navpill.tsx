"use client";

import { Presence } from "@ark-ui/react";

import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";
import { Input } from "src/theme/components/Input";

import { CloseAction } from "../../Action/Close";
import { AdminAction } from "../Anchors/Admin";
import { ComposeAction } from "../Anchors/Compose";
import { DashboardAction } from "../Anchors/Dashboard";
import { HomeAction } from "../Anchors/Home";
import { LinksAction } from "../Anchors/Links";
import { LoginAction } from "../Anchors/Login";
import { LogoutAction } from "../Anchors/Logout";
import { NotificationsAction } from "../Anchors/Notifications";
import { SettingsAction } from "../Anchors/Settings";

import { css } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";

import { Menu } from "./components/Menu";
import { SearchResults } from "./components/SearchResults";
import { useNavpill } from "./useNavpill";

export function Navpill() {
  const {
    isExpanded,
    onExpand,
    onClose,
    account,
    searchQuery,
    onSearch,
    searchResults,
  } = useNavpill();
  return (
    <Toolpill onClickOutside={onClose}>
      <Presence present={isExpanded} className={css({ w: "full" })}>
        <styled.div w="full">
          <HStack w="full" justify="space-between">
            {account ? (
              <>
                <HStack>
                  <HomeAction />
                  <NotificationsAction />
                  <LinksAction />
                  <LogoutAction />
                </HStack>
                <HStack>
                  {account.admin && (
                    <>
                      <AdminAction />
                    </>
                  )}
                  <SettingsAction />
                </HStack>
              </>
            ) : (
              <LoginAction />
            )}
          </HStack>

          {searchResults.length ? (
            <SearchResults results={searchResults} />
          ) : (
            <Menu />
          )}
        </styled.div>
      </Presence>

      {account ? (
        <HStack gap="4" w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <ProfilePill profileReference={account} showHandle={false} />

              <Input
                border="none"
                size="sm"
                placeholder="Search disabled..."
                disabled
                value={searchQuery}
                onChange={onSearch}
              />
              <CloseAction onClick={onClose} />
            </>
          ) : (
            <>
              <ProfilePill profileReference={account} showHandle={false} />
              <HomeAction />
              <ComposeAction />
              <NotificationsAction />
              <DashboardAction onClick={onExpand} />
            </>
          )}
        </HStack>
      ) : (
        <HStack gap="4" w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <HomeAction />
              <Input
                border="none"
                placeholder="Search disabled..."
                disabled
                value={searchQuery}
                onChange={onSearch}
              />
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
