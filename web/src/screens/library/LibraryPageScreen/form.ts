import { z } from "zod";

import {
  PropertySchema,
  PropertySchemaList,
  PropertyType,
} from "src/api/openapi-schema";

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

// NOTE: Yes, SchemaSchema is intentional...
export const FormNodeChildPropertySchemaSchema = z.object({
  fid: z.string(),
  name: z.string(),
  sort: z.string(),
  type: z.enum(PropertyTypes),
});
export type FormNodeChildPropertySchema = z.infer<
  typeof FormNodeChildPropertySchemaSchema
>;
