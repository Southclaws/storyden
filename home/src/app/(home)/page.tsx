import {
  CareChapter,
  CommunityInvitation,
  ConversationChapter,
  KnowledgeChapter,
  OwnershipChapter,
  StoryIntroduction,
} from "./sections/CommunityStory";
import { HeroSection } from "./sections/HeroSection";
import { type HomeStats, MilspecSection } from "./sections/MilspecSection";
import { ScreenshotSection } from "./sections/ScreenshotSection";

export default async function Home() {
  const stats = await getStats();
  return (
    <div>
      <HeroSection />
      <div className="marketing-story">
        <ScreenshotSection />
        <StoryIntroduction />
        <ConversationChapter />
        <KnowledgeChapter />
        <CareChapter />
        <OwnershipChapter />
        <MilspecSection stats={stats} />
        <CommunityInvitation />
      </div>
    </div>
  );
}

async function getStats(): Promise<HomeStats> {
  // Unavailable figures stay explicitly unavailable; no historical guesses.
  const stats: HomeStats = {
    stars: null,
    commits: null,
    contributors: null,
    loc: null,
    apis: null,
  };
  try {
    const response = await fetch(
      "https://api.github.com/repos/Southclaws/storyden",
      {
        headers: { Accept: "application/vnd.github+json" },
        next: { revalidate: 86400 },
        signal: AbortSignal.timeout(2500),
      },
    );
    if (response.ok) {
      const data = await response.json();
      if (
        typeof data.stargazers_count === "number" &&
        Number.isFinite(data.stargazers_count)
      ) {
        stats.stars = data.stargazers_count;
      }
    }
  } catch (error) {
    console.warn(
      "Homepage repository statistics unavailable:",
      error instanceof Error ? error.message : "request failed",
    );
  }
  return stats;
}
