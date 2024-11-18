import { z } from "zod";

import { AuthMode } from "@/api/openapi-schema";

export const AuthenticationModeTable: { [K in AuthMode]: string } = {
  handle: "handle",
  email: "email",
  phone: "phone",
};

const AuthenticationModes = Object.keys(
  AuthenticationModeTable,
) as unknown as readonly [AuthMode, ...AuthMode[]];

export const AuthenticationModeSchema: z.ZodType<AuthMode> =
  z.enum(AuthenticationModes);
export type AuthenticationMode = z.infer<typeof AuthenticationModeSchema>;

export type AuthenticationModeDetail = {
  value: AuthMode;
  name: string;
  description: string;
};

export const AuthenticationModeDetails: Record<
  AuthMode,
  AuthenticationModeDetail
> = {
  [AuthMode.handle]: {
    value: AuthMode.handle,
    name: "Username + Password",
    description:
      "The simplest authentication mode, members register and log in with their username and password. This mode does not require an email client to be configured and is suitable for smaller or invite-only communities, however it leaves the community vulnerable to spam and abuse as members cannot be verified by email. Members will not be able to reset their password by themselves and administrators will not have a contact method for members.",
  },
  [AuthMode.email]: {
    value: AuthMode.email,
    name: "Email",
    description:
      "Email authentication is the most flexible and common mode. Members register and log in with their email address and password. This mode requires an email client to be configured. Members can reset their password by themselves and administrators have a contact method for members.",
  },
  [AuthMode.phone]: {
    value: AuthMode.phone,
    name: "Phone",
    description:
      "This mode enforces members to register and log in with their just their phone number via one-time verification codes instead of passwords. It's recommended for communities that are mobile-first and requires an SMS client to be configured.",
  },
};

export const AuthenticationModeList = Object.values(AuthenticationModeDetails);
