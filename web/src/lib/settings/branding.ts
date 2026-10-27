import * as z from "@/lib/validation";

export const BrandingSettingsSchema = z.object({
  banner: z.optional(
    z.object({
      asset_id: z.string(),
      crop: z.object({
        left: z.number().check(z.minimum(0), z.maximum(1)),
        top: z.number().check(z.minimum(0), z.maximum(1)),
        width: z.number().check(z.positive(), z.maximum(1)),
        height: z.number().check(z.positive(), z.maximum(1)),
      }),
    }),
  ),
});
