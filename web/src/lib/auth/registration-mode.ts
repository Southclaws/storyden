import { z } from "zod";

import { RegistrationMode } from "@/api/openapi-schema";

export const RegistrationModeTable: { [K in RegistrationMode]: string } = {
  public: "public",
  invitation: "invitation",
  disabled: "disabled",
};

const RegistrationModes = Object.keys(
  RegistrationModeTable,
) as unknown as readonly [RegistrationMode, ...RegistrationMode[]];

export const RegistrationModeSchema = z.enum(RegistrationModes);

export type RegistrationModeDetail = {
  value: RegistrationMode;
  name: string;
  description: string;
};

export const RegistrationModeDetails: Record<
  RegistrationMode,
  RegistrationModeDetail
> = {
  [RegistrationMode.public]: {
    value: RegistrationMode.public,
    name: "Public",
    description:
      "Anyone can create an account using the selected authentication mode.",
  },
  [RegistrationMode.invitation]: {
    value: RegistrationMode.invitation,
    name: "Invitation",
    description:
      "New members need an invitation link. Normal public registration links are hidden.",
  },
  [RegistrationMode.disabled]: {
    value: RegistrationMode.disabled,
    name: "Disabled",
    description:
      "Self-serve registration is closed. New accounts can only be provisioned through the API.",
  },
};

export const RegistrationModeList = Object.values(RegistrationModeDetails);
