import Image from "next/image";

// These original passages remain independently placeable in the marketing story.
export function Permanence() {
  return (
    <div className="permanence-section">
      <p>
        In a tangled web of fleeting <strong>feeds</strong> and forgotten
        threads, Storyden makes <strong>what matters</strong> stay discoverable,
        readable, and shareable. <strong>On your terms</strong>.
      </p>
      <p>
        Own your <strong>data</strong>. Build your own <strong>corner</strong>{" "}
        of the web.
      </p>
      <p className="permanence-signoff">
        Collective memory. Without the noise.
      </p>
    </div>
  );
}

export function Organise() {
  return (
    <div className="original-diagram">
      <h3>Conversations flow. Ideas grow.</h3>
      <p>
        A good answer can lead to a shared reading list, and a shared reading
        list can start the next conversation.
      </p>
      <Image
        src="/organise.png"
        alt="A member asks for design book recommendations; a reply links to a community-curated Design Books Directory containing a book about Dieter Rams."
        width={1320}
        height={864}
        sizes="(max-width: 760px) 90vw, 510px"
      />
    </div>
  );
}

export function Gardens() {
  return (
    <div className="gardens-section">
      <Image
        src="/square-tree-smol.webp"
        alt=""
        width={1024}
        height={1024}
        sizes="(max-width: 760px) 140px, 220px"
      />
      <div>
        <h3>
          Knowledge gardens,
          <br />
          not content farms.
        </h3>
        <p>
          From fan <strong>clubs</strong> to research groups, digital zines to
          esports <strong>teams</strong>. For anyone who cares about ideas and
          the <strong>people</strong> behind them.
        </p>
      </div>
    </div>
  );
}

export function Accessible() {
  return (
    <section className="accessible-section">
      <Image src="/accessibility.png" width={64} height={64} alt="" />
      <div>
        <h2>Optimised for humans.</h2>
        <p>
          A place for people to share, discover, and build knowledge together.
        </p>
      </div>
    </section>
  );
}

export function Curate() {
  return (
    <div className="original-diagram">
      <h3>Curate, with a little help.</h3>
      <p>
        <span className="optional-label">Optional AI</span> Bring in a link and
        get a starting point for its title, tags, and summary.
      </p>
      <Image
        src="/curate.png"
        alt="Importing a link into the Library: the website is fetched, tags are assigned, and a Library page is created."
        width={1320}
        height={864}
        sizes="(max-width: 760px) 90vw, 510px"
      />
    </div>
  );
}

// Retained as a complete composition for future repositioning.
export function CollectiveMemorySection() {
  return (
    <div className="marketing-story">
      <Permanence />
      <Organise />
      <Gardens />
      <Curate />
    </div>
  );
}
