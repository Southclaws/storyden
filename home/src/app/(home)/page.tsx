import type { Metadata } from "next";

import { getHomeStats } from "@/lib/home-stats";

import {
  CareChapter,
  CommunityInvitation,
  ConversationChapter,
  KnowledgeChapter,
  OwnershipChapter,
  StoryIntroduction,
} from "./sections/CommunityStory";
import { HeroSection } from "./sections/HeroSection";
import { MilspecSection } from "./sections/MilspecSection";
import { ScreenshotSection } from "./sections/ScreenshotSection";

export const metadata: Metadata = {
  alternates: { canonical: "https://www.storyden.org/" },
};

export default async function Home() {
  const stats = await getHomeStats();
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
