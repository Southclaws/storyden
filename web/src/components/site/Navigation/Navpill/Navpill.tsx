"use client";

import {
  Admin,
  Bell,
  Close,
  Create,
  Dashboard,
  Home,
  Login,
  Logout,
  Settings,
} from "src/components/site/Action/Action";
import { ProfilePill } from "src/components/site/ProfilePill/ProfilePill";
import { Toolpill } from "src/components/site/Toolpill/Toolpill";
import { HStack, Input, SlideFade } from "src/theme/components";

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
      <SlideFade
        in={isExpanded}
        style={{
          maxHeight: "100%",
          width: "100%",
          display: isExpanded ? "flex" : "none",
          flexDirection: "column",
        }}
      >
        <HStack justify="space-between">
          {account ? (
            <>
              <HStack>
                <Home />
                <Bell />
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
      </SlideFade>

      {account ? (
        <HStack gap={4} w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <ProfilePill profileReference={account} showHandle={false} />

              <Input
                variant="outline"
                border="none"
                size="sm"
                placeholder="Search disabled..."
                isDisabled
                value={searchQuery}
                onChange={onSearch}
              />
              <Close onClick={onClose} />
            </>
          ) : (
            <>
              <ProfilePill profileReference={account} showHandle={false} />
              <Home />
              <Create />
              <Bell />
              <Dashboard onClick={onExpand} />
            </>
          )}
        </HStack>
      ) : (
        <HStack gap={4} w="full" justifyContent="space-between">
          {isExpanded ? (
            <>
              <Home />
              <Input
                variant="outline"
                border="none"
                placeholder="Search disabled..."
                isDisabled
                value={searchQuery}
                onChange={onSearch}
              />
              <Close onClick={onClose} />
            </>
          ) : (
            <>
              <Home />
              <Login />
              <Dashboard onClick={onExpand} />
            </>
          )}
        </HStack>
      )}
    </Toolpill>
  );
}
