import { z } from "zod";

export const UsernameSchema = z
  .string()
  .min(1, "Please enter a username.")
  .max(30, "Maximum length is 30 characters.")
  .toLowerCase()
  .regex(
    /^[a-z0-9_-]+$/g,
    "Username can only contain latin letters, numbers, dashes and underscores.",
  );

export const PasswordSchema = z
  .string()
  .min(8, "Password must be at least 8 characters.");

export const ExistingPasswordSchema = z
  .string()
  .min(1, "Please enter your current password.");
