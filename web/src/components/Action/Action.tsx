import { Link, LinkProps } from "@chakra-ui/next-js";
import { IconButton, IconButtonProps } from "@chakra-ui/react";
import {
  AdjustmentsHorizontalIcon,
  Bars3Icon,
  BellIcon,
  HomeIcon,
  PlusIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { LoginIcon } from "../graphics/LoginIcon";
import { SpeechPlusIcon } from "../graphics/SpeechPlusIcon";
import { MouseEvent, MouseEventHandler, useCallback } from "react";

function useClickHandler(onClick: MouseEventHandler | undefined) {
  // This allows us to progressively enhance features on the application by
  // treating important buttons as links to fallback pages. For example, there
  // may be a button that triggers the opening of a modal dialogue but if the
  // user has JavaScript disabled due to device constraints or privacy reasons,
  // the functionality must also be implemented by a normal page.
  return useCallback(
    (e: MouseEvent) => {
      if (onClick) {
        e.preventDefault();
        return onClick?.(e);
      }
    },
    [onClick]
  );
}

export function Action({ children, onClick, ...props }: LinkProps) {
  const handleClick = useClickHandler(onClick);
  return (
    <Link
      width={8}
      height={8}
      display="flex"
      justifyContent="center"
      alignItems="center"
      borderRadius="full"
      p={1}
      _hover={{ bgColor: "blackAlpha.50" }}
      onClick={handleClick}
      {...props}
    >
      {children}
    </Link>
  );
}

export function ActionButton({ children, ...props }: IconButtonProps) {
  return (
    <IconButton
      width={8}
      height={8}
      display="flex"
      justifyContent="center"
      alignItems="center"
      borderRadius="full"
      p={1}
      bgColor="transparent"
      _hover={{ bgColor: "blackAlpha.50" }}
      {...props}
    >
      {children}
    </IconButton>
  );
}

// A few actions have default page destinations (partly for consistency and also
// for accessibility and no-JS modes) so this just redefines `href` as optional.
type WithOptionalURL = Omit<LinkProps, "href"> & {
  href?: string | undefined;
};

export function Bell({ href = "/notifications", ...props }: WithOptionalURL) {
  return (
    <Action href={href} {...props}>
      <BellIcon width="1.25em" />
    </Action>
  );
}

export function Home({ href = "/", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="Home" {...props}>
      <HomeIcon width="1.25em" />
    </Action>
  );
}

export function Login({ href = "/auth", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="Sign up or Log in" {...props}>
      <LoginIcon width="1.5em" />
    </Action>
  );
}

export function Create({ href = "/new", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="New thread" {...props}>
      <SpeechPlusIcon width="1.25em" />
    </Action>
  );
}

export function Dashboard({ href = "/dashboard", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="Dashboard" {...props}>
      <Bars3Icon width="1.25em" />
    </Action>
  );
}

export function Settings({ href = "/settings", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="Settings" {...props}>
      <AdjustmentsHorizontalIcon width="1.25em" />
    </Action>
  );
}

type WithOptionalARIALabel = Omit<IconButtonProps, "aria-label"> & {
  "aria-label"?: string | undefined;
};

export function Close({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  return (
    <ActionButton
      size="sm"
      title="Close"
      aria-label={al ?? "close"}
      {...props}
      icon={<XMarkIcon width="1.4em" />}
    />
  );
}

export function Add({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  return (
    <ActionButton
      size="sm"
      title="Add"
      aria-label={al ?? "add"}
      {...props}
      icon={<PlusIcon width="1.4em" />}
    />
  );
}
