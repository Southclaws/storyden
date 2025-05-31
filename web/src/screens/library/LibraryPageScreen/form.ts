import { z } from "zod";

import { PropertyType } from "src/api/openapi-schema";

import { CoverImageSchema, NodeMetadataSchema } from "@/lib/library/metadata";

const CoverImageFormSchema = z.union([
  CoverImageSchema,
  z.object({
    asset_id: z.string(),
  }),
]);

const PropertyTypes = Object.keys(PropertyType) as unknown as readonly [
  PropertyType,
  ...PropertyType[],
];

export const FormNodePropertySchema = z.object({
  fid: z.string().optional(),
  name: z.string(),
  type: z.enum(PropertyTypes),
  sort: z.string(),
  value: z.string(),
});
export type FormNodeProperty = z.infer<typeof FormNodePropertySchema>;

export const FormSchema = z.object({
  name: z.string().min(1, "Please enter a name."),
  slug: z.string().optional(),
  properties: z.array(FormNodePropertySchema),
  tags: z.string().array().optional(),
  link: z.preprocess((v) => {
    if (typeof v === "string" && v === "") {
      return undefined;
    }

    return v;
  }, z.string().url("Invalid URL").optional()),
  coverImage: CoverImageFormSchema.optional(),
  content: z.string().optional(),
  meta: NodeMetadataSchema,
});
export type Form = z.infer<typeof FormSchema>;
