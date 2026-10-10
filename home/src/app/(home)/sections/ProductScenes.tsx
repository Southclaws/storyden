import {
  ArrowDown,
  ArrowRight,
  BookOpen,
  Check,
  Clock3,
  Ellipsis,
  Heart,
  Library,
  Reply,
  Route,
  Search,
  TriangleAlert,
} from "lucide-react";
import Image from "next/image";

import { ConnectionMotion } from "./ConnectionMotion";

// Simplified from web's ThreadScreen, Byline, and Reply/CardBox.
// Keep this static illustration in the same reading order as the application.
export function DiscussionScene() {
  return (
    <figure className="discussion-example">
      <div className="product-scene thread-scene">
        <div className="thread-preview-main">
          <div className="thread-preview-byline">
            <span className="scene-avatar avatar-moss" aria-hidden="true">
              J
            </span>
            <span>@jules</span>
            <span>· 2h</span>
            <span className="thread-preview-category">Work in progress</span>
          </div>
          <h3>The little harbour is finally playable</h3>
          <p className="thread-preview-body">
            I got the boats moving! The lighting still needs work. How would you
            make the docks feel a little more lived-in?
          </p>
          <div className="thread-preview-stats">
            <span>
              <Heart size={15} aria-hidden="true" /> 8 likes
            </span>
            <span>
              <Image
                src="/illustrations/discussion.svg"
                width={15}
                height={15}
                alt=""
              />{" "}
              4 replies
            </span>
          </div>
          <div className="thread-preview-reply">
            <div className="thread-preview-reply-header">
              <div className="thread-preview-byline">
                <span className="scene-avatar avatar-clay" aria-hidden="true">
                  M
                </span>
                <span>@mika</span>
                <span>· 1h</span>
              </div>
              <span className="thread-preview-reply-actions" aria-hidden="true">
                <Reply size={15} />
                <Ellipsis size={17} />
              </span>
            </div>
            <p>
              Love the scale of this. The{" "}
              <span className="thread-preview-reference">
                environment art field guide
              </span>{" "}
              in our Library has a whole section on telling stories with props.
            </p>
            <div className="thread-preview-emoji">
              <span role="img" aria-label="Raised hands">
                🙌
              </span>{" "}
              2
            </div>
          </div>
          <div className="thread-preview-reply">
            <div className="thread-preview-reply-header">
              <div className="thread-preview-byline">
                <span className="scene-avatar avatar-moss" aria-hidden="true">
                  J
                </span>
                <span>@jules</span>
                <span>· 45min</span>
              </div>
              <span className="thread-preview-reply-actions" aria-hidden="true">
                <Reply size={15} />
                <Ellipsis size={17} />
              </span>
            </div>
            <p>
              Exactly what I needed. Some coiled ropes and a lantern by the
              boats next!
            </p>
          </div>
        </div>
      </div>
    </figure>
  );
}

export function LibrarySequence() {
  return (
    <figure className="library-example">
      <div className="library-sequence product-scene">
        <div className="library-step">
          <div className="scene-person">
            <span className="scene-avatar avatar-clay">A</span>
            <strong>Alex shared a find</strong>
          </div>
          <h3>Good design book recommendations?</h3>
          <p>“This one changed how I think about everyday objects.”</p>
        </div>
        <div className="sequence-connector">
          <ArrowDown size={22} aria-hidden="true" />
          <span>Keep it in the Library</span>
        </div>
        <div className="library-step">
          <div className="scene-toolbar">
            <span>
              <Library size={16} aria-hidden="true" /> Design reading room
            </span>
            <span className="scene-meta">Library</span>
          </div>
          <div className="directory-row">
            <BookOpen size={18} aria-hidden="true" />
            <strong>The Design of Everyday Things</strong>
            <span className="scene-tag">Book</span>
          </div>
          <dl className="scene-properties">
            <div>
              <dt>Author</dt>
              <dd>Don Norman</dd>
            </div>
            <div>
              <dt>Topic</dt>
              <dd>Human-centred design</dd>
            </div>
          </dl>
        </div>
        <div className="sequence-connector">
          <ArrowDown size={22} aria-hidden="true" />
          <span>Useful to the next person too!</span>
        </div>
        <div className="library-step library-found">
          <Search size={19} aria-hidden="true" />
          <div>
            <strong>Start here: learning design</strong>
            <p>A shared reading list, ready when someone asks.</p>
          </div>
          <Check size={18} aria-hidden="true" />
        </div>
      </div>
    </figure>
  );
}

export function TrailScene() {
  return (
    <figure className="trail-example">
      <div className="product-scene trail-scene">
        <div className="scene-toolbar">
          <span>
            <Route size={16} aria-hidden="true" /> Sunday photography club
          </span>
          <span className="scene-tag">Optional AI</span>
        </div>
        <div className="scene-content">
          <h3>A little help before the weekend</h3>
          <p className="scene-meta">
            A Trail runs a configured Robot on your schedule.
          </p>
          <ol className="trail-timeline">
            <li>
              <Clock3 size={18} aria-hidden="true" />
              <div>
                <strong>Friday, 16:00</strong>
                <p>Review this week’s Library contributions.</p>
              </div>
              <span className="scene-tag">Scheduled</span>
            </li>
            <li>
              <Check size={18} aria-hidden="true" />
              <div>
                <strong>Review completed</strong>
                <p>Three new location guides added this week.</p>
              </div>
              <span className="scene-tag tag-success">Completed</span>
            </li>
            <li>
              <TriangleAlert size={18} aria-hidden="true" />
              <div>
                <strong>One for the organisers</strong>
                <p>A suggested change needs a decision.</p>
              </div>
              <span className="scene-tag tag-attention">Needs attention</span>
            </li>
          </ol>
        </div>
      </div>
    </figure>
  );
}

export function ConnectionDiagram() {
  return (
    <figure className="connection-example">
      <ConnectionMotion>
        <div
          className="connection-map"
          role="img"
          aria-label="Your Storyden community connects to custom websites and tools through its API, OAuth, and plugins. A separate branch leads to optional Robots, MCP, and AI connections."
        >
          <div className="connection-source">
            <div className="connection-source-mark">
              <Image
                src="/brand/logomark_newspaper_600.png"
                alt=""
                width={56}
                height={56}
              />
            </div>
            <div className="connection-source-copy">
              <strong>Your Storyden</strong>
              <span>Your members, content, and rules.</span>
            </div>
          </div>
          <ConnectionLines />
          <ConnectionLines mobile />
          <div className="connection-branches">
            <div className="connection-branch">
              <strong>Make it your own</strong>
              <p>Custom sites & tools</p>
              <span>API · OAuth · Plugins</span>
            </div>
            <div className="connection-branch">
              <strong>Add help, if you want it</strong>
              <p>Robots & connected tools</p>
              <span>MCP · Optional AI</span>
            </div>
          </div>
        </div>
      </ConnectionMotion>
    </figure>
  );
}

function ConnectionLines({ mobile = false }: { mobile?: boolean }) {
  const paths = mobile
    ? [
        "M200 0V24",
        "M200 24C200 56 100 40 100 80",
        "M200 24C200 56 300 40 300 80",
      ]
    : [
        "M-72 128H24",
        "M24 128C76 128 40 48 104 48H172",
        "M24 128C76 128 40 208 104 208H172",
      ];

  return (
    <svg
      className={`connection-lines ${mobile ? "connection-lines-mobile" : "connection-lines-desktop"}`}
      width={mobile ? 400 : 180}
      height={mobile ? 96 : 256}
      viewBox={mobile ? "0 0 400 96" : "0 0 180 256"}
      preserveAspectRatio={mobile ? "none" : "xMidYMid meet"}
      fill="none"
      aria-hidden="true"
    >
      {paths.map((path, index) => (
        <path
          key={path}
          d={path}
          className={
            index === 2
              ? "connection-route connection-route-optional"
              : "connection-route"
          }
        />
      ))}
      <path
        d={mobile ? "M96 84L100 88L104 84" : "M176 44L180 48L176 52"}
        className="connection-arrow"
      />
      <path
        d={mobile ? "M296 84L300 88L304 84" : "M176 204L180 208L176 212"}
        className="connection-arrow connection-arrow-optional"
      />
    </svg>
  );
}

// Based on the rendered reports page and web/screens/reports/manager/ReportCard.
export function ModerationExample() {
  return (
    <figure
      className="moderation-example product-scene"
      aria-label="A submitted report about a discussion becoming personal, awaiting moderator review"
    >
      <div className="moderation-report-header">
        <div>
          <span className="moderation-status">Submitted</span>
          <span className="moderation-kind">Thread</span>
        </div>
        <span className="scene-meta">20min ago</span>
      </div>
      <blockquote>
        “The feedback here has turned into personal remarks. Could a moderator
        take a look?”
      </blockquote>
      <div className="moderation-people">
        <div>
          <span>Reporter:</span>
          <span className="scene-avatar avatar-clay" aria-hidden="true">
            A
          </span>
          @alex
        </div>
        <div>
          <span>Handler:</span>
          <span className="moderation-unassigned">Unassigned</span>
        </div>
      </div>
      <div className="moderation-subject">
        <div className="moderation-subject-heading">
          <h3>Feedback on the harbour scene</h3>
          <span className="scene-meta">1h ago</span>
        </div>
        <div className="moderation-author">
          <span className="scene-avatar avatar-moss" aria-hidden="true">
            J
          </span>
          <span>@jules</span>
        </div>
      </div>
      <div className="moderation-actions">
        <span className="moderation-acknowledge">Acknowledge</span>
        <span className="moderation-resolve">Resolve</span>
      </div>
    </figure>
  );
}
