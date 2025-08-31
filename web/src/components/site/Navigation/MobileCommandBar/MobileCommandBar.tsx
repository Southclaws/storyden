"use client";

import { CommandDock } from "@/components/site/CommandDock/CommandDock";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { SiteIcon } from "@/components/ui/icons/Site";
import { WStack } from "@/styled-system/jsx";

import { Search } from "../../../search/Search/Search";
import { CloseAction } from "../../Action/Close";
import { AccountMenu } from "../AccountMenu/AccountMenu";
import { ComposeAnchor } from "../Anchors/Compose";
import { HomeAnchor } from "../Anchors/Home";
import { LibraryAnchor } from "../Anchors/Library";
import { LoginAnchor } from "../Anchors/Login";
import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { useMobileCommandBar } from "./useMobileCommandBar";

export function MobileCommandBar() {
  const { isExpanded, onExpand, onClose, account } = useMobileCommandBar();

  return (
    <CommandDock
      isOpen={isExpanded}
      onClickOutside={onClose}
      render={() => {
        return <ContentNavigationList />;
      }}
    >
      <WStack alignItems="center">
        {isExpanded ? (
          <>
            {account ? (
              <AccountMenu account={account} size="sm" />
            ) : (
              <SiteIcon borderRadius="md" w="8" h="8" />
            )}
            <Search />
            <CloseAction onClick={onClose} size="sm" />
          </>
        ) : (
          <>
            {account ? (
              <AccountMenu account={account} size="sm" />
            ) : (
              <SiteIcon borderRadius="md" w="8" h="8" />
            )}
            <HomeAnchor hideLabel size="sm" />
            {account ? <ComposeAnchor hideLabel size="sm" /> : <LoginAnchor />}
            <LibraryAnchor hideLabel size="sm" />
            <ExpandTrigger onClick={onExpand} />
          </>
        )}
      </WStack>
    </CommandDock>
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
