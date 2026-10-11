import { execFileSync } from "node:child_process";
import { readFileSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { parse } from "yaml";

import { getHomeStats } from "../src/lib/home-stats";
import snapshot from "../src/lib/home-stats.snapshot.json";

async function main() {
  const root = fileURLToPath(new URL("../../", import.meta.url));
  const schema = parse(
    readFileSync(new URL("../../api/openapi.yaml", import.meta.url), "utf8"),
  );
  const methods = new Set([
    "get",
    "put",
    "post",
    "delete",
    "options",
    "head",
    "patch",
    "trace",
  ]);
  snapshot.apis = Object.values(schema.paths).reduce<number>(
    (total, path) =>
      total +
      Object.keys(path as object).filter((key) => methods.has(key)).length,
    0,
  );

  try {
    const output = execFileSync(
      "tokei",
      [
        "--output",
        "json",
        "--exclude",
        "home/src/lib/home-stats.snapshot.json",
        ".",
      ],
      {
        cwd: root,
        encoding: "utf8",
        maxBuffer: 20 * 1024 * 1024,
      },
    );
    snapshot.loc = JSON.parse(output).Total.code;
    snapshot.locMeasuredAt = new Date().toISOString().slice(0, 10);
  } catch (error) {
    // Vercel doesn't provide tokei. Keep the last actual measurement there.
    if ((error as NodeJS.ErrnoException).code !== "ENOENT") throw error;
    console.warn(
      `tokei unavailable; retaining LOC measured ${snapshot.locMeasuredAt}`,
    );
  }

  if (process.argv.includes("--github")) {
    const stats = await getHomeStats({ strict: true });
    snapshot.stars = stats.stars;
    snapshot.commits = stats.commits;
    snapshot.contributors = stats.contributors;
    snapshot.githubMeasuredAt = new Date().toISOString().slice(0, 10);
  }

  writeFileSync(
    new URL("../src/lib/home-stats.snapshot.json", import.meta.url),
    `${JSON.stringify(snapshot, null, 2)}\n`,
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
