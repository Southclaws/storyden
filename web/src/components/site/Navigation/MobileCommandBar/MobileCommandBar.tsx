"use client";

import { useSession } from "@/auth";
import type { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { MenuIcon } from "@/components/ui/icons/Menu";
import { SiteIcon } from "@/components/ui/icons/Site";

import { Search } from "../../../search/Search/Search";
import { CloseAction } from "../../Action/Close";
import { AccountMenu } from "../AccountMenu/AccountMenu";
import { ComposeAnchor } from "../Anchors/Compose";
import { HomeAnchor } from "../Anchors/Home";
import { LibraryAnchor } from "../Anchors/Library";
import { LoginAnchor } from "../Anchors/Login";

type Props = {
  canRegister?: boolean;
  isOpen: boolean;
  onClose: () => void;
  onOpen: () => void;
};

export function MobileCommandBar({
  canRegister,
  isOpen,
  onClose,
  onOpen,
}: Props) {
  const account = useSession();

  return (
    <div
      aria-label="Quick navigation"
      className="navigation__mobile-commandbar"
      role="toolbar"
    >
      <div className="navigation__mobile-commandbar-inner">
        <div className="navigation__mobile-commandbar-items">
          {isOpen ? (
            <>
              {account ? (
                <AccountMenu account={account} size="sm" />
              ) : (
                <span className="navigation__mobile-site-icon">
                  <SiteIcon />
                </span>
              )}
              <Search />
              <CloseAction
                aria-label="Close navigation menu"
                onClick={onClose}
                size="sm"
                title="Close navigation menu"
                type="button"
              />
            </>
          ) : (
            <>
              {account ? (
                <AccountMenu account={account} size="sm" />
              ) : (
                <span className="navigation__mobile-site-icon">
                  <SiteIcon />
                </span>
              )}
              <HomeAnchor hideLabel size="sm" />
              {account ? (
                <ComposeAnchor hideLabel size="sm" />
              ) : (
                canRegister && <LoginAnchor />
              )}
              <LibraryAnchor hideLabel size="sm" />
              <ExpandTrigger isOpen={isOpen} onClick={onOpen} />
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function ExpandTrigger({
  isOpen,
  ...props
}: ButtonProps & { isOpen: boolean }) {
  return (
    <IconButton
      aria-controls="navigation__leftbar"
      aria-expanded={isOpen}
      aria-label="Open navigation menu"
      title="Open navigation menu"
      type="button"
      variant="ghost"
      size="sm"
      {...props}
    >
      <MenuIcon />
    </IconButton>
  );
}
