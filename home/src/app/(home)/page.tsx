import { Box } from "@/styled-system/jsx";

import { CollectiveMemorySection } from "./sections/CollectiveMemorySection";
import { HeroSection } from "./sections/HeroSection";
import { MilspecSection, type HomeStats } from "./sections/MilspecSection";
import { ScreenshotSection } from "./sections/ScreenshotSection";

export default async function Home() {
  const stats = await getStats();

  return (
    <Box>
      <HeroSection />
      <ScreenshotSection />
      <CollectiveMemorySection />
      <MilspecSection stats={stats} />
    </Box>
  );
}

async function getStats(): Promise<HomeStats> {
  const defaults: HomeStats = {
    // 2025-07-05
    stars: 279,
    commits: 2702,
    contributors: 17,
    loc: 314217, // tokei --output json | jq .Total.code
    apis: 140, // rg operationId ./api/openapi.yaml -c
  };

  try {
    const repo = "Southclaws/storyden";

    const headers = {
      Accept: "application/vnd.github+json",
      // Uncomment below and add a token if you hit rate limits:
      // Authorization: `Bearer ${process.env.GITHUB_TOKEN}`,
    };

    const [repoRes, contributorsRes, commitsRes] = await Promise.all([
      fetch(`https://api.github.com/repos/${repo}`, { headers }),
      fetch(`https://api.github.com/repos/${repo}/contributors?per_page=100`, {
        headers,
      }),
      fetch(`https://api.github.com/repos/${repo}/commits?per_page=1`, {
        headers,
      }),
    ]);

    if (!repoRes.ok || !contributorsRes.ok || !commitsRes.ok) {
      return defaults;
    }

    const repoData = (await repoRes.json()) as {
      stargazers_count?: number;
      open_issues_count?: number;
    };
    const contributors = (await contributorsRes.json()) as Array<{
      contributions?: number;
    }>;

    const openIssues = Number(repoData.open_issues_count ?? 0);
    const contributorCommits = contributors.reduce(
      (acc, contributor) => acc + Number(contributor.contributions ?? 0),
      0
    );

    return {
      ...defaults,
      stars: Number(repoData.stargazers_count ?? defaults.stars),
      commits: openIssues + contributorCommits,
      contributors: contributors.length,
    };
  } catch {
    return defaults;
  }
}
