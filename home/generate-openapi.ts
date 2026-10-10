import { readFileSync, writeFileSync } from "fs";
import { generateFiles } from "fumadocs-openapi";
import { parse } from "yaml";

import { openapi } from "@/lib/openapi";

// Generate the API documentation files
void generateFiles({
  input: openapi,
  output: "./content/docs/api",
  per: "operation",
  groupBy: "tag",
  includeDescription: true,
  frontmatter: (title, description) => ({
    title,
    description,
    full: false,
  }),
});

// Parse the OpenAPI spec to get tag ordering
const openapiSpec = parse(readFileSync("../api/openapi.yaml", "utf-8"));
const tags = openapiSpec.tags?.map((tag: { name: string }) => tag.name) || [];

// Generate meta.json with the tag order
const meta = {
  root: true,
  title: "API Reference",
  description: "The Storyden RESTful API.",
  pages: ["index", ...tags],
};

writeFileSync(
  "./content/docs/api/meta.json",
  JSON.stringify(meta, null, 2) + "\n",
);
