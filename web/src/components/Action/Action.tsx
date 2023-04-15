import { Link, LinkProps } from "@chakra-ui/next-js";
import { HomeIcon } from "@heroicons/react/24/outline";
import { BellIcon } from "@heroicons/react/24/solid";

export function Action({ children, ...props }: LinkProps) {
  return <Link {...props}>{children}</Link>;
}

// A few actions have default page destinations (partly for consistency and also
// for accessibility and no-JS modes) so this just redefines `href` as optional.
type WithOptionalURL = Omit<LinkProps, "href"> & {
  href?: string | undefined;
};

export function Bell(props: LinkProps) {
  return (
    <Action {...props}>
      <BellIcon />
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
      <HomeIcon width="1.5em" />
    </Action>
  );
}
