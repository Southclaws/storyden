import { Link } from "src/theme/components/Link";

import { StyleProps } from "@/styled-system/types";

export function LoginAction(props: StyleProps) {
  return (
    <Link href="/login" kind="secondary" size="sm" {...props}>
      Login
    </Link>
  );
}

export function RegisterAction(props: StyleProps) {
  return (
    <Link href="/register" kind="primary" size="sm" {...props}>
      Register
    </Link>
  );
}
