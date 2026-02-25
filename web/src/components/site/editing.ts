import { z } from "zod";

export const EditingSchema = z.preprocess(
  (value) => {
    if (typeof value === "string" && value === "") {
      return undefined;
    }

    return value;
  },
  z.enum(["settings", "feed"]),
);
export type Editing = z.infer<typeof EditingSchema>;
