import type { ReactNode } from "react";

import { addressRules } from "./address-rules";
import { clientIPRules } from "./client-ip-rules";

export type ClientIPContext = {
  mode: "remote_addr" | "single_header" | "xff_trusted_proxies";
  trustedProxyCIDRs: string[];
  forwardedFor?: string;
  resolvedAddress?: string;
};

export type DoctorContext = {
  api: URL;
  web: URL;
  browser?: URL;
  clientIP?: ClientIPContext;
};

type DoctorAdvice = {
  title: ReactNode;
  explanation: ReactNode;
  remedy: ReactNode;
};

export type DoctorRule = {
  code: string;
  category: "addresses" | "authentication" | "routing" | "networking";
  severity: "warning" | "error";
  check: (context: DoctorContext) => DoctorAdvice | undefined;
};

export const doctorRules: readonly DoctorRule[] = [
  ...addressRules,
  ...clientIPRules,
];

export type DoctorFinding = DoctorAdvice &
  Pick<DoctorRule, "code" | "category" | "severity">;

export function diagnoseInstallation(
  apiAddress: string,
  webAddress: string,
  browserOrigin?: string,
  clientIP?: ClientIPContext,
): DoctorFinding[] {
  const context: DoctorContext = {
    api: new URL(apiAddress),
    web: new URL(webAddress),
    browser: browserOrigin ? new URL(browserOrigin) : undefined,
    clientIP,
  };

  return doctorRules.flatMap((rule) => {
    const advice = rule.check(context);

    if (!advice) {
      return [];
    }

    return [
      {
        code: rule.code,
        category: rule.category,
        severity: rule.severity,
        ...advice,
      },
    ];
  });
}
