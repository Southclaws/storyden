"use client";

import Link from "next/link";
import { useState } from "react";

import { AccountCommonProps, AuthMode, Permission } from "@/api/openapi-schema";
import { Admonition } from "@/components/ui/admonition";
import { Settings } from "@/lib/settings/settings";
import { Box } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

type Props = {
  session: AccountCommonProps | undefined;
  settings: Settings;
};

export function VerificationBanner({ session, settings }: Props) {
  const [visible, setVisible] = useState(true);

  if (!session) {
    return null;
  }

  const isAdmin = hasPermission(session, Permission.ADMINISTRATOR);

  const shouldShow =
    settings.authentication_mode === AuthMode.email &&
    session.verified_status === "none" &&
    !isAdmin;

  if (!shouldShow) {
    return null;
  }

  return (
    <Box mb="4">
      <Admonition
        value={visible}
        kind="failure"
        title="Email Verification Required"
        onChange={setVisible}
      >
        <p>
          Please{" "}
          <Link
            href="/settings?tab=email"
            style={{ textDecoration: "underline" }}
          >
            verify your email in settings
          </Link>{" "}
          to participate in this community.
        </p>
      </Admonition>
    </Box>
  );
}
