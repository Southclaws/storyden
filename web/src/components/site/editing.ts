import * as z from "@/lib/validation";

export const EditingSchema = z.pipe(
  z.transform((value) => {
    if (typeof value === "string" && value === "") {
      return undefined;
    }

    if (value === "feed") {
      return "site";
    }

    return value;
  }),
  z.optional(z.enum(["settings", "site"])),
);
export type Editing = z.infer<typeof EditingSchema>;
