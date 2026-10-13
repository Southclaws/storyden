import * as z from "@/lib/validation";
import { isSlug } from "@/utils/slugify";

export const UsernameSchema = z.string().check(
  z.minLength(1, "Please enter a username."),
  z.maxLength(30, "Maximum length is 30 characters."),
  z.toLowerCase(),
  z.refine(
    (val) => isSlug(val),
    "Username must be lowercase letters, numbers, hyphens, and underscores only.",
  ),
);

export const PasswordSchema = z
  .string()
  .check(z.minLength(8, "Password must be at least 8 characters."));

export const ExistingPasswordSchema = z
  .string()
  .check(z.minLength(1, "Please enter your current password."));

export const UsernameOrEmailSchema = z.union([
  UsernameSchema,
  z.string().check(z.email("Please enter a valid email.")),
]);
