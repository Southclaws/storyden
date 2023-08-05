import { z } from "zod";

export const ParamSchema = z.object({
  handle: z.string().optional(),
  collection: z.string().optional(),
});
