import { ArrowLeftIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { Link, LinkProps } from "src/theme/components/Link";

export function BackAction({
  children,
  ...props
}: PropsWithChildren<LinkProps>) {
  return <Link {...props}>{children ?? <ArrowLeftIcon width="1.4em" />}</Link>;
}
