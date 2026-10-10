import { ArrowRight, ExternalLink } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

import {
  Curate,
  Gardens,
  Organise,
  Permanence,
} from "./CollectiveMemorySection";
import {
  ConnectionDiagram,
  DiscussionScene,
  LibrarySequence,
  ModerationExample,
  TrailScene,
} from "./ProductScenes";

export function StoryIntroduction() {
  return (
    <section className="story-introduction" aria-label="The Storyden manifesto">
      <div className="story-container">
        <Permanence />
      </div>
    </section>
  );
}

export function ConversationChapter() {
  return (
    <section
      className="story-chapter story-conversation"
      id="conversations"
      aria-labelledby="conversation-title"
    >
      <div className="story-container">
        <div className="story-split conversation-layout">
          <div className="story-copy">
            <h1 id="conversation-title">
              Home for your people.
              <br />
              <span>And for what they know.</span>
            </h1>
            <p className="story-lead">Community you can come back to.</p>
            <p>
              A place for conversations to unfold, ideas to take root, and
              people to stay connected at their own pace.
            </p>
            <Link
              className="story-text-link"
              href="/docs/introduction/discussion"
            >
              Room for a proper conversation{" "}
              <ArrowRight size={17} aria-hidden="true" />
            </Link>
          </div>
          <DiscussionScene />
        </div>
        <Gardens />
      </div>
    </section>
  );
}

export function KnowledgeChapter() {
  return (
    <section
      className="story-chapter story-knowledge"
      id="library"
      aria-labelledby="library-title"
    >
      <div className="story-container">
        <div className="story-split knowledge-layout">
          <div className="story-copy">
            <h2 id="library-title">
              Make the good finds
              <br />
              part of the furniture.
            </h2>
            <p className="story-lead">
              Someone shared it.
              <br />
              Now everyone can build on it.
            </p>
            <p>
              A reading list for your design circle. A wiki for your game. A
              directory of places your club loves. The Library gives shared
              knowledge a home beside the conversations that made it.
            </p>
            <p>
              Organise pages together, add useful properties, and link people to
              something that keeps getting better.
            </p>
            <Link className="story-text-link" href="/docs/introduction/library">
              Meet your community’s Library{" "}
              <ArrowRight size={17} aria-hidden="true" />
            </Link>
          </div>
          <LibrarySequence />
        </div>
        <figure className="library-screenshot">
          <Image
            src="/docs/introduction/library/directory.png"
            alt="Storyden’s Design Library: a page tree beside a searchable directory of design resources, with descriptions and contribution status."
            width={2130}
            height={1700}
            sizes="(max-width: 760px) 100vw, 1100px"
          />
        </figure>
        <div className="knowledge-details">
          <Organise />
          <Curate />
        </div>
        <div className="optional-note">
          <span className="optional-label">Optional AI</span>
          <p>
            <strong>A helping hand with the filing.</strong> Suggest a title or
            tags, import a summary from a link, or ask a configured Robot to
            help find and organise material. You can also do it all yourself.
          </p>
          <Link href="/docs/introduction/autofill" className="story-text-link">
            See how assistance works <ArrowRight size={17} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  );
}

export function CareChapter() {
  return (
    <section
      className="story-chapter story-care"
      id="looking-after"
      aria-labelledby="care-title"
    >
      <div className="story-container">
        <div className="care-introduction">
          <div className="story-copy">
            <h2 id="care-title">
              A shared space,
              <br />
              looked after together.
            </h2>
            <p className="story-lead">
              Running a community is a labour of love.
            </p>
            <p>
              A community’s character comes from the people who care for it.
              Thoughtful tools for moderation and shared spaces help keep that
              character intact as the community grows.
            </p>
            <Link
              className="story-text-link"
              href="/docs/introduction/moderation"
            >
              Look after your community{" "}
              <ArrowRight size={17} aria-hidden="true" />
            </Link>
          </div>
          <ModerationExample />
        </div>
        <div className="care-assistance">
          <div className="story-copy">
            <h3>
              Some of the routine
              <br />
              can look after itself.
            </h3>
            <span className="optional-label">Optional AI</span>
            <p>
              Configure a Robot to help with a focused job, then use a Trail to
              run it on a schedule or in response to an event.
            </p>
            <p>
              Review recent contributions. Surface things that need attention.
              Keep the run history close, and leave decisions that need a person
              with your organisers.
            </p>
            <p>
              You choose the model provider, tools, and permissions. AI
              assistance requires configuration; it’s never a condition of
              joining in.
            </p>
            <Link href="/docs/introduction/robots" className="story-text-link">
              Explore connected assistance{" "}
              <ArrowRight size={17} aria-hidden="true" />
            </Link>
          </div>
          <TrailScene />
        </div>
      </div>
    </section>
  );
}

export function OwnershipChapter() {
  return (
    <section
      className="story-chapter story-ownership"
      id="your-community"
      aria-labelledby="ownership-title"
    >
      <div className="story-container">
        <div className="ownership-introduction">
          <h2 id="ownership-title">
            Your community.
            <br />
            <span>Your terms.</span>
          </h2>
          <div className="story-copy">
            <p>
              Put your name on the door. Set your own rules. Run Storyden on
              your own infrastructure, with your community’s identity front and
              centre. Big social tech has no say.
            </p>
            <p>
              Start with the built-in defaults. Add integrations, change the
              frontend, or connect your own tools as your ambitions grow.
            </p>
            <Link
              href="/docs/introduction/vps/docker"
              className="story-text-link"
            >
              Make yourself at home <ArrowRight size={17} aria-hidden="true" />
            </Link>
          </div>
        </div>
        <ConnectionDiagram />
        <div className="ownership-links">
          <Link href="/docs/reference/settings#brand-settings">
            Your look and feel <ArrowRight size={16} aria-hidden="true" />
          </Link>
          <Link href="/docs/extending">
            Your integrations <ArrowRight size={16} aria-hidden="true" />
          </Link>
          <a href="https://github.com/Southclaws/storyden">
            Everyone&apos;s source code <ExternalLink size={15} aria-hidden="true" />
          </a>
        </div>
      </div>
    </section>
  );
}

export function CommunityInvitation() {
  return (
    <section
      className="community-invitation"
      aria-labelledby="invitation-title"
    >
      <Image
        src="/square-campfire.webp"
        width={1024}
        height={1024}
        alt=""
        sizes="(max-width: 760px) 180px, 300px"
      />
      <div>
        <h2 id="invitation-title">
          Come on in.
          <br />
          Make yourself at home.
        </h2>
        <p>
          Meet the people behind Storyden in our community.
          <br />
          Then imagine what yours could become.
        </p>
        <div className="story-actions">
          <Link className="story-button" href="https://makeroom.club">
            Explore our community{" "}
            <ArrowRight size={18} aria-hidden="true" />
          </Link>
          <Link className="story-text-link" href="/docs/introduction">
            Start your own <ArrowRight size={17} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  );
}
