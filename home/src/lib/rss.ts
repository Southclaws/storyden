import { Feed } from "feed";
import { blog } from "@/lib/source";
import { readFileSync } from "fs";
import { join } from "path";
import { marked } from "marked";

const baseUrl = "https://www.storyden.org";

export async function getRSS() {
  const feed = new Feed({
    title: "Storyden Blog",
    id: `${baseUrl}/blog`,
    link: `${baseUrl}/blog`,
    description:
      "Storyden is a platform for building communities. A modern take on oldschool bulletin board forums. Designed to be the community platform for the next era of internet culture.",
    language: "en",

    image: `${baseUrl}/banner.png`,
    favicon: `${baseUrl}/icon.png`,
    copyright: "Barnaby Keene",
  });

  const pages = blog.getPages();

  const sorted = pages.sort((a, b) => {
    return new Date(b.data.date!).getTime() - new Date(a.data.date!).getTime();
  });

  async function getPostContent(
    page: (typeof pages)[number]
  ): Promise<string | undefined> {
    try {
      const filePath = join(process.cwd(), "content/blog", page.path);
      const content = readFileSync(filePath, "utf-8");

      const bodyContent = content.replace(/^---[\s\S]*?---\n/, "");

      const html = await marked.parse(bodyContent);

      return html;
    } catch (error) {
      console.warn(`Failed to load content for ${page.url}:`, error);
      return undefined;
    }
  }

  for (const page of sorted) {
    const content = await getPostContent(page);

    feed.addItem({
      id: page.url,
      title: page.data.title,
      description: page.data.description,
      content: content,
      link: `${baseUrl}${page.url}`,
      date: new Date(page.data.date!),

      author: [
        {
          name: "Barnaby Keene",
        },
      ],
    });
  }

  return feed.rss2();
}
