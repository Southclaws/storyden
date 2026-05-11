"use client";

import { LinkButton } from "@/components/ui/link-button";
import { useI18n } from "@/i18n/provider";
import { JsxStyleProps } from "@/styled-system/types";

export function LoginAnchor(props: JsxStyleProps) {
  const { t } = useI18n();

  return (
    <LinkButton href="/login" variant="ghost" size="sm" {...props}>
      {t("Login")}
    </LinkButton>
  );
}

export function RegisterAnchor(props: JsxStyleProps) {
  const { t } = useI18n();

  return (
    <LinkButton href="/register" size="sm" {...props}>
      {t("Register")}
    </LinkButton>
  );
}
