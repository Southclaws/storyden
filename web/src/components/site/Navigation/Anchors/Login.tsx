import { LinkButton } from "@/components/ui/link-button";
import { JsxStyleProps } from "@/styled-system/types";

export function LoginAnchor(props: JsxStyleProps) {
  return (
    <LinkButton href="/login" variant="ghost" size="sm" {...props}>
      Login
    </LinkButton>
  );
}

export function RegisterAnchor(props: JsxStyleProps) {
  return (
    <LinkButton href="/register" size="sm" {...props}>
      Register
    </LinkButton>
  );
}
