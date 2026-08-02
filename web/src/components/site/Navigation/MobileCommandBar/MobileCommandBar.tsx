"use client";

import { useSession } from "@/auth";
import type { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CloseIcon } from "@/components/ui/icons/Close";
import { MenuIcon } from "@/components/ui/icons/Menu";

import { AccountMenu } from "../AccountMenu/AccountMenu";
import { HomeAnchor } from "../Anchors/Home";
import { LoginAnchor } from "../Anchors/Login";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onOpen: () => void;
};

export function MobileCommandBar({ isOpen, onClose, onOpen }: Props) {
  const account = useSession();

  return (
    <div
      aria-label="Primary navigation"
      className="navigation__mobile-commandbar navigation__surface"
      role="toolbar"
    >
      <NavigationDrawerTrigger
        isOpen={isOpen}
        onClick={isOpen ? onClose : onOpen}
      />

      <div className="navigation__mobile-commandbar-end">
        <HomeAnchor hideLabel size="sm" />
        {account ? (
          <AccountMenu account={account} size="sm" />
        ) : (
          <LoginAnchor />
        )}
      </div>
    </div>
  );
}

function NavigationDrawerTrigger({
  isOpen,
  ...props
}: ButtonProps & { isOpen: boolean }) {
  const label = isOpen ? "Close navigation menu" : "Open navigation menu";

  return (
    <IconButton
      aria-controls="navigation__leftbar"
      aria-expanded={isOpen}
      aria-label={label}
      type="button"
      variant="ghost"
      size="sm"
      {...props}
    >
      {isOpen ? <CloseIcon /> : <MenuIcon />}
    </IconButton>
  );
}
