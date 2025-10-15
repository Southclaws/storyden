import { generateFiles } from "fumadocs-openapi";
import { openapi } from "@/lib/openapi";

void generateFiles({
  input: openapi,
  output: "./content/docs/api",
  per: "tag",
  includeDescription: true,
  frontmatter: (title, description) => ({
    title,
    description,
    full: false,
  }),
});
