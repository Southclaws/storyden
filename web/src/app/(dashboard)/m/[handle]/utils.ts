import * as z from "@/lib/validation";

export const ParamSchema = z.object({
  handle: z.optional(z.string()),
  collection: z.optional(z.string()),
});
