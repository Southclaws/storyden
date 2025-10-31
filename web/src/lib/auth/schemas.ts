import { z } from "zod";

import { isSlug } from "@/utils/slugify";

export const UsernameSchema = z
  .string()
  .min(1, "Please enter a username.")
  .max(30, "Maximum length is 30 characters.")
  .toLowerCase()
  .refine(
    (val) => isSlug(val),
    "Username must be lowercase letters, numbers, hyphens, and underscores only.",
  );

export const PasswordSchema = z
  .string()
  .min(8, "Password must be at least 8 characters.");

export const ExistingPasswordSchema = z
  .string()
  .min(1, "Please enter your current password.");

export const UsernameOrEmailSchema = z.union([
  UsernameSchema,
  z.string().email("Please enter a valid email."),
]);
