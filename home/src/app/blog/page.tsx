import {blog} from "@/lib/source";
import {VStack} from "@/styled-system/jsx";

export default function Page() {
  const posts = [...blog.getPages()].sort(
    (a, b) =>
      new Date(b.data.date ?? b.file.name).getTime() -
      new Date(a.data.date ?? a.file.name).getTime()
  );

  return (
    <VStack>
      <h1>Blog</h1>
      <ul>
        {posts.map((post) => (
          <li key={post.file.name}>
            <a href={post.url}>
              <h1>{post.data.title}</h1>
              <p>{post.data.description}</p>
            </a>
          </li>
        ))}
      </ul>

      <pre>wip...</pre>
    </VStack>
  );
}
