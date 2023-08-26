"use client";

import { HStack, Input, SlideFade } from "@chakra-ui/react";

import {
  Bell,
  Close,
  Create,
  Dashboard,
  Home,
  Login,
  Logout,
  Settings,
} from "src/components/Action/Action";
import { ProfileReference } from "src/components/ProfileReference/ProfileReference";
import { Toolpill } from "src/components/Toolpill/Toolpill";

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
              <Settings />
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
              <ProfileReference profileReference={account} showHandle={false} />

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
              <ProfileReference profileReference={account} showHandle={false} />
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
