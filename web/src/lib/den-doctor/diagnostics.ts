import type { ReactNode } from "react";

import { addressRules } from "./address-rules";

export type DoctorContext = {
  api: URL;
  web: URL;
  browser?: URL;
};

type DoctorAdvice = {
  title: ReactNode;
  explanation: ReactNode;
  remedy: ReactNode;
};

export type DoctorRule = {
  code: string;
  category: "addresses" | "authentication" | "routing";
  severity: "warning" | "error";
  check: (context: DoctorContext) => DoctorAdvice | undefined;
};

export const doctorRules: readonly DoctorRule[] = [...addressRules];

export type DoctorFinding = DoctorAdvice &
  Pick<DoctorRule, "code" | "category" | "severity">;

export function diagnoseInstallation(
  apiAddress: string,
  webAddress: string,
  browserOrigin?: string,
): DoctorFinding[] {
  const context: DoctorContext = {
    api: new URL(apiAddress),
    web: new URL(webAddress),
    browser: browserOrigin ? new URL(browserOrigin) : undefined,
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
