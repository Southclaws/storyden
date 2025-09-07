import {
  defineDocs,
  defineConfig,
  defineCollections,
  frontmatterSchema,
} from "fumadocs-mdx/config";
import { z } from "zod";

export const docs = defineDocs({
  dir: "content/docs",
});

export const blog = defineCollections({
  type: "doc",
  dir: "content/blog",
  async: true,
  schema: frontmatterSchema.extend({
    // vercel builds differently for some dumb unknown (the usual)
    // so hack this to accept both types and transform... thanks vercel.
    date: z
      .union([z.string(), z.date()])
      .transform((val) => (val instanceof Date ? val : new Date(val))),
  }),
});

export default defineConfig({
  mdxOptions: {
    // MDX options
  },
});
