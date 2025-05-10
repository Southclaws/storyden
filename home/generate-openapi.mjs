import { generateFiles } from "fumadocs-openapi";

void generateFiles({
  input: [
    "https://raw.githubusercontent.com/Southclaws/storyden/refs/heads/main/api/openapi.yaml",
  ],
  output: "./content/docs/api",
  per: "tag",
});
