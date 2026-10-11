import snapshot from "./home-stats.snapshot.json";

const repository = "https://api.github.com/repos/Southclaws/storyden";

async function github(path: string) {
  const response = await fetch(`${repository}${path}`, {
    headers: { Accept: "application/vnd.github+json" },
    next: { revalidate: 86400 },
    signal: AbortSignal.timeout(5000),
  });
  if (!response.ok) throw new Error(`GitHub returned ${response.status}`);
  return response;
}

async function count(path: string): Promise<number> {
  // One item per page makes the last page number the exact item count.
  const response = await github(`${path}?per_page=1`);
  const items: unknown = await response.json();
  if (!Array.isArray(items)) throw new Error("Expected a GitHub item list");
  const links = response.headers.get("link") ?? "";
  const last = links.match(/<([^>]+)>;\s*rel="last"/);
  if (last) {
    const total = Number(new URL(last[1]).searchParams.get("page"));
    if (Number.isSafeInteger(total) && total > 0) return total;
    throw new Error("Invalid GitHub last-page count");
  }
  if (links.includes('rel="next"')) {
    throw new Error("GitHub pagination has no total");
  }
  return items.length;
}

export async function getHomeStats({ strict = false } = {}) {
  const stats = {
    stars: snapshot.stars,
    commits: snapshot.commits,
    contributors: snapshot.contributors,
    loc: snapshot.loc,
    apis: snapshot.apis,
  };
  const results = await Promise.allSettled([
    github("").then(async (response) => {
      const data = await response.json();
      if (
        !Number.isSafeInteger(data.stargazers_count) ||
        data.stargazers_count < 0
      ) {
        throw new Error("Invalid GitHub star count");
      }
      return data.stargazers_count as number;
    }),
    count("/commits"),
    count("/contributors"),
  ]);
  const fields = ["stars", "commits", "contributors"] as const;
  results.forEach((result, index) => {
    const field = fields[index];
    if (result.status === "fulfilled") {
      stats[field] = result.value;
    } else {
      if (strict) throw result.reason;
      console.warn(`Homepage ${field}: using measured snapshot`, result.reason);
    }
  });
  return stats;
}
