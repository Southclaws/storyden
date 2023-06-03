import { Link, LinkProps } from "@chakra-ui/next-js";

import {
  IconButton,
  IconButtonProps,
  Link as DumbLink,
  forwardRef,
} from "@chakra-ui/react";
import {
  AdjustmentsHorizontalIcon,
  ArrowLeftIcon,
  Bars3Icon,
  BellIcon,
  CloudArrowUpIcon,
  EllipsisHorizontalIcon,
  HomeIcon,
  PlusIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { MouseEvent, MouseEventHandler, useCallback } from "react";
import { LoginIcon } from "../graphics/LoginIcon";
import { LogoutIcon } from "../graphics/LogoutIcon";
import { SpeechPlusIcon } from "../graphics/SpeechPlusIcon";
import { SendIcon } from "../graphics/SendIcon";
import { useCommands } from "@remirror/react";

const actionStyles = {
  width: 8,
  height: 8,
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  borderRadius: "full",
  p: 1,
  _hover: { bgColor: "blackAlpha.50" },
};

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
    <Link onClick={handleClick} {...actionStyles} {...props}>
      {children}
    </Link>
  );
}

export const ActionButton = forwardRef<IconButtonProps, "button">(
  ({ children, ...props }, ref) => {
    return (
      <IconButton
        ref={ref}
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
);

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

export function Logout({ href = "/logout", ...props }: WithOptionalURL) {
  // We need to use a regular link here (Chakra styled anchor tag) because the
  // anchor tag provided by Next.js is too clever for logouts! Because the Link
  // component from Next.js pre-loads pages when the user hovers, this results
  // in unexpected logouts just from hovering over the logout button, not ideal!
  return (
    <DumbLink
      href={href}
      title="Log out of your session"
      {...actionStyles}
      {...props}
    >
      <LogoutIcon width="1.5em" />
    </DumbLink>
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

export function Back({
  href,
  "aria-label": al,
  ...props
}: WithOptionalARIALabel & WithOptionalURL) {
  if (href)
    return (
      <Action
        href={href}
        size="sm"
        title="Back"
        aria-label={al ?? "back"}
        {...props}
      >
        <ArrowLeftIcon width="1.4em" />
      </Action>
    );

  return (
    <ActionButton
      size="sm"
      title="Back"
      aria-label={al ?? "back"}
      {...props}
      icon={<ArrowLeftIcon width="1.4em" />}
    />
  );
}

export function Send({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  return (
    <ActionButton
      size="sm"
      title="Send"
      aria-label={al ?? "send"}
      {...props}
      icon={<SendIcon width="1.4em" />}
    />
  );
}

export function Save({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  return (
    <ActionButton
      size="sm"
      title="Save"
      aria-label={al ?? "save"}
      {...props}
      icon={<CloudArrowUpIcon width="1.4em" />}
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

// NOTE: this one is forward-ref'd because it's used as a chakra Menu button.
// https://chakra-ui.com/docs/components/menu#customizing-the-button
export const More = forwardRef(
  ({ "aria-label": al, ...props }: WithOptionalARIALabel, ref) => {
    return (
      <ActionButton
        ref={ref}
        size="sm"
        title="More"
        aria-label={al ?? "more"}
        {...props}
        icon={<EllipsisHorizontalIcon width="1.4em" />}
      />
    );
  }
);

export function Bold({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  const { toggleBold, focus } = useCommands();
  return (
    <ActionButton
      onClick={() => {
        toggleBold();
        focus();
      }}
      aria-label={al ?? "Bold"}
      icon={
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M13.259 17.625H7.47778V6.375H12.7903C13.4166 6.37504 14.0299 6.55435 14.5576 6.89174C15.0853 7.22914 15.5054 7.71052 15.7683 8.27903C16.0312 8.84754 16.1259 9.47942 16.0412 10.1C15.9565 10.7206 15.6959 11.304 15.2903 11.7812C15.8201 12.205 16.2056 12.7825 16.3936 13.4344C16.5816 14.0862 16.563 14.7803 16.3403 15.4211C16.1175 16.0619 15.7016 16.6179 15.1498 17.0126C14.598 17.4073 13.9374 17.6213 13.259 17.625ZM9.35278 15.75H13.2465C13.4312 15.75 13.6141 15.7136 13.7847 15.643C13.9553 15.5723 14.1103 15.4687 14.2409 15.3381C14.3715 15.2075 14.4751 15.0525 14.5457 14.8819C14.6164 14.7113 14.6528 14.5284 14.6528 14.3438C14.6528 14.1591 14.6164 13.9762 14.5457 13.8056C14.4751 13.635 14.3715 13.48 14.2409 13.3494C14.1103 13.2188 13.9553 13.1152 13.7847 13.0445C13.6141 12.9739 13.4312 12.9375 13.2465 12.9375H9.35278V15.75ZM9.35278 11.0625H12.7903C12.975 11.0625 13.1578 11.0261 13.3284 10.9555C13.499 10.8848 13.6541 10.7812 13.7847 10.6506C13.9152 10.52 14.0188 10.365 14.0895 10.1944C14.1602 10.0238 14.1965 9.84092 14.1965 9.65625C14.1965 9.47158 14.1602 9.28872 14.0895 9.1181C14.0188 8.94749 13.9152 8.79246 13.7847 8.66188C13.6541 8.5313 13.499 8.42772 13.3284 8.35704C13.1578 8.28637 12.975 8.25 12.7903 8.25H9.35278V11.0625Z"
            fill="#212529"
          />
        </svg>
      }
      {...props}
    />
  );
}

export function Italic({ "aria-label": al, ...props }: WithOptionalARIALabel) {
  const { toggleItalic, focus } = useCommands();
  return (
    <ActionButton
      onClick={() => {
        toggleItalic();
        focus();
      }}
      aria-label={al ?? "Italic"}
      icon={
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M17.625 7.625V6.375H9.5V7.625H12.7125L9.98125 16.375H6.375V17.625H14.5V16.375H11.2875L14.0187 7.625H17.625Z"
            fill="#212529"
          />
        </svg>
      }
      {...props}
    />
  );
}

export function Underline({
  "aria-label": al,
  ...props
}: WithOptionalARIALabel) {
  const { toggleUnderline, focus } = useCommands();
  return (
    <ActionButton
      onClick={() => {
        toggleUnderline();
        focus();
      }}
      aria-label={al ?? "Italic"}
      icon={
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M4.5 17.9375H19.5V19.1875H4.5V17.9375ZM12 16.0625C10.8397 16.0625 9.72688 15.6016 8.90641 14.7811C8.08594 13.9606 7.625 12.8478 7.625 11.6875V4.8125H8.875V11.6875C8.875 12.5163 9.20424 13.3112 9.79029 13.8972C10.3763 14.4833 11.1712 14.8125 12 14.8125C12.8288 14.8125 13.6237 14.4833 14.2097 13.8972C14.7958 13.3112 15.125 12.5163 15.125 11.6875V4.8125H16.375V11.6875C16.375 12.8478 15.9141 13.9606 15.0936 14.7811C14.2731 15.6016 13.1603 16.0625 12 16.0625V16.0625Z"
            fill="#212529"
          />
        </svg>
      }
      {...props}
    />
  );
}
