import { Link } from "src/theme/components/Link";

import { StyleProps } from "@/styled-system/types";

export function LoginAction(props: StyleProps) {
  return (
    <Link href="/login" variant="ghost" size="sm" {...props}>
      Login
    </Link>
  );
}

export function RegisterAction(props: StyleProps) {
  return (
    <Link href="/register" size="sm" {...props}>
      Register
    </Link>
  );
}
