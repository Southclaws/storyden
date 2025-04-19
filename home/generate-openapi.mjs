import { generateFiles } from "fumadocs-openapi";

void generateFiles({
  input: ["../api/openapi.yaml"],
  output: "./content/docs/api",
  groupBy: "tag",
});
