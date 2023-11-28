"use client";

import { Presence } from "@ark-ui/react";

import {
  Admin,
  Close,
  Login,
  Logout,
  Settings,
} from "src/components/site/Action/Action";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";
import { Input } from "src/theme/components/Input";

import { ComposeAction } from "../../Action/Compose";
import { DashboardAction } from "../../Action/Dashboard";
import { HomeAction } from "../../Action/Home";
import { LinksAction } from "../../Action/Links";
import { NotificationsAction } from "../../Action/Notifications";

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
      <Presence present={isExpanded}>
        <styled.div w="full">
          <HStack w="full" justify="space-between">
            {account ? (
              <>
                <HStack>
                  <HomeAction />
                  <NotificationsAction />
                  <LinksAction />
                  <Logout />
                </HStack>
                <HStack>
                  {account.admin && (
                    <>
                      <Admin />
                    </>
                  )}
                  <Settings />
                </HStack>
              </>
            ) : (
              <Login />
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
              <Close onClick={onClose} />
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
              <Close onClick={onClose} />
            </>
          ) : (
            <>
              <HomeAction />
              <Login />
              <DashboardAction onClick={onExpand} />
            </>
          )}
        </HStack>
      )}
    </Toolpill>
  );
}
