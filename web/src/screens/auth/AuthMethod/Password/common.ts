import * as z from "zod";

export const FormSchema = z.object({
  identifier: z.string(),
  token: z.string(),
});
export type Form = z.infer<typeof FormSchema>;
