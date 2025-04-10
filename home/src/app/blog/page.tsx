import { blog } from "@/lib/source";
import { Card, styled } from "@/styled-system/jsx";
import { vstack } from "@/styled-system/patterns";
import { formatDate, formatDistanceToNow } from "date-fns";
import { DocsPage } from "fumadocs-ui/page";

export default function Page() {
  const posts = [...blog.getPages()].sort(
    (a, b) =>
      new Date(b.data.date ?? b.file.name).getTime() -
      new Date(a.data.date ?? a.file.name).getTime()
  );

  return (
    <DocsPage>
      <h1>Storyden Blog</h1>
      <styled.ul className={vstack()} gap="2">
        {posts.map((post) => (
          <Card key={post.url}>
            <li key={post.file.name}>
              <a href={post.url}>
                <h2>{post.data.title}</h2>
                <p>{post.data.description}</p>
                <styled.time
                  color="slate.500"
                  title={formatDate(post.data.date!, "yyyy-MM-dd")}
                >
                  {formatDistanceToNow(post.data.date!, { addSuffix: true })}
                </styled.time>
              </a>
            </li>
          </Card>
        ))}
      </styled.ul>
    </DocsPage>
  );
}
