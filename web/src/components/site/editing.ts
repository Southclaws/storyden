import { z } from "zod";

export const EditingSchema = z.preprocess(
  (value) => {
    if (typeof value === "string" && value === "") {
      return undefined;
    }

    if (value === "feed") {
      return "site";
    }

    return value;
  },
  z.enum(["settings", "site"]).optional(),
);
export type Editing = z.infer<typeof EditingSchema>;
