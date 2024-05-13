/**
 * Generated by orval v6.28.2 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
 * OpenAPI spec version: 1
 */

/**
 * The phone number request payload.
 */
export interface PhoneRequestCodeProps {
  /** The desired username to link to the phone number. */
  identifier: string;
  /** The phone number to receive the one-time code on. */
  phone_number: string;
}
