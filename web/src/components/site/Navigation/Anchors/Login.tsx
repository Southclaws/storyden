import { LinkButton } from "@/components/ui/link-button";
import { StyleProps } from "@/styled-system/types";

export function LoginAction(props: StyleProps) {
  return (
    <LinkButton href="/login" variant="ghost" size="sm" {...props}>
      Login
    </LinkButton>
  );
}

export function RegisterAction(props: StyleProps) {
  return (
    <LinkButton href="/register" size="sm" {...props}>
      Register
    </LinkButton>
  );
}
