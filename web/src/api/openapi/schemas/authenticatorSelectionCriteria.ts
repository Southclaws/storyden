/**
 * Generated by orval v6.30.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */
import type { AuthenticatorAttachment } from "./authenticatorAttachment";
import type { ResidentKeyRequirement } from "./residentKeyRequirement";
import type { UserVerificationRequirement } from "./userVerificationRequirement";

/**
 * https://www.w3.org/TR/webauthn-2/#dictdef-authenticatorselectioncriteria

 */
export interface AuthenticatorSelectionCriteria {
  authenticatorAttachment: AuthenticatorAttachment;
  requireResidentKey?: boolean;
  residentKey: ResidentKeyRequirement;
  userVerification?: UserVerificationRequirement;
}
