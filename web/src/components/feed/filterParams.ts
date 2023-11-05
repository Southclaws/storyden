import { z } from "zod";

export const FilterParamsSchema = z.object({
  category: z.string().optional(),
});
export type FilterParams = z.infer<typeof FilterParamsSchema>;
