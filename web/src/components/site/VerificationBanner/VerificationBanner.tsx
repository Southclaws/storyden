"use client";

import Link from "next/link";
import { useState } from "react";

import { AccountCommonProps, AuthMode } from "@/api/openapi-schema";
import { Admonition } from "@/components/ui/admonition";
import { useI18n } from "@/i18n/provider";
import { Settings } from "@/lib/settings/settings";
import { Box } from "@/styled-system/jsx";

type Props = {
  session: AccountCommonProps | undefined;
  settings: Settings;
};

export function VerificationBanner({ session, settings }: Props) {
  const [visible, setVisible] = useState(true);
  const { t } = useI18n();

  if (!session) {
    return null;
  }

  const shouldShow =
    settings.authentication_mode === AuthMode.email &&
    session.verified_status === "none";

  if (!shouldShow) {
    return null;
  }

  return (
    <Box mb="4">
      <Admonition
        value={visible}
        kind="failure"
        title={t("Email Verification Required")}
        onChange={setVisible}
      >
        <p>
          {t("Please")}{" "}
          <Link
            href="/settings?tab=email"
            style={{ textDecoration: "underline" }}
          >
            {t("verify your email in settings")}
          </Link>{" "}
          {t("to participate in this community.")}
        </p>
      </Admonition>
    </Box>
  );
}
