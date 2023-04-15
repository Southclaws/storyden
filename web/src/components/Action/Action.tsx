import { Link, LinkProps } from "@chakra-ui/next-js";
import { BellIcon } from "@heroicons/react/24/solid";

export function Action({ children, ...props }: LinkProps) {
  return <Link {...props}>{children}</Link>;
}

export function Bell(props: LinkProps) {
  return (
    <Action {...props}>
      <BellIcon />
    </Action>
  );
}
