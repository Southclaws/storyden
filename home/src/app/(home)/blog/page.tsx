import { formatDate } from "date-fns";
import Image from "next/image";
import Link from "next/link";

import { RSSButtons } from "@/components/RSSButtons";
import { blog } from "@/lib/source";

export default function Page() {
  const posts = [...blog.getPages()].sort(
    (a, b) => new Date(b.data.date).getTime() - new Date(a.data.date).getTime(),
  );
  return (
    <div className="blog-index">
      <header className="blog-heading">
        <div className="blog-heading-copy">
          <h1>From the workshop.</h1>
          <p>
            Notes on building Storyden, the communities that inspire it, and a
            more personal corner of the internet.
          </p>
          <RSSButtons />
        </div>
        <Image
          className="blog-artwork"
          alt="An illustrated mountain peak above a pine-lined valley in warm evening light."
          src="/wide-mountain.webp"
          width={2048}
          height={1024}
          sizes="(max-width: 700px) calc(100vw - 40px), (max-width: 1000px) 280px, 420px"
        />
      </header>
      <div className="blog-post-list">
        {posts.map((post) => (
          <article className="blog-post" key={post.url}>
            <time dateTime={formatDate(post.data.date, "yyyy-MM-dd")}>
              {formatDate(post.data.date, "d MMMM yyyy")}
            </time>
            <h2>
              <Link href={post.url}>{post.data.title}</Link>
            </h2>
            <p>{post.data.description}</p>
          </article>
        ))}
      </div>
    </div>
  );
}
