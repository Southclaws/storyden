import { Link, LinkProps } from "@chakra-ui/next-js";
import { HomeIcon } from "@heroicons/react/24/outline";
import { BellIcon } from "@heroicons/react/24/outline";
import { LoginIcon } from "../graphics/LoginIcon";
import { SpeechPlusIcon } from "../graphics/SpeechPlusIcon";

export function Action({ children, ...props }: LinkProps) {
  return (
    <Link
      borderRadius="full"
      p={1}
      _hover={{ bgColor: "blackAlpha.50" }}
      {...props}
    >
      {children}
    </Link>
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

export function Menu(props: LinkProps) {
  return (
    <Action {...props}>
      <BellIcon />
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
    <Action href={href} title="Home" {...props}>
      <LoginIcon width="1.5em" />
    </Action>
  );
}

export function Create({ href = "/new", ...props }: WithOptionalURL) {
  return (
    <Action href={href} title="Home" {...props}>
      <SpeechPlusIcon width="1.25em" />
    </Action>
  );
}
