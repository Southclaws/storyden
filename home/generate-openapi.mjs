import { generateFiles } from "fumadocs-openapi";

void generateFiles({
  input: [
    "https://raw.githubusercontent.com/Southclaws/storyden/refs/heads/main/api/openapi.yaml",
  ],
  output: "./content/docs/api",

  // NOTE: Currently broken so don't try to re-generate docs!
  // Waiting for a fumadocs fix. Right now it will generate both folders and
  // files with duplicate docs content when grouped by tags.
  groupBy: "tags",
  per: "tag",
});
