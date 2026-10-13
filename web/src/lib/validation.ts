import { config } from "zod/mini";
import en from "zod/v4/locales/en.js";

// Mini has no default locale. Preserve the existing English validation messages.
config(en());

export * from "zod/mini";
