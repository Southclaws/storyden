import { Feed } from "feed";
import { blog } from "@/lib/source";

const baseUrl = "https://www.storyden.org";

export function getRSS() {
  const feed = new Feed({
    title: "Storyden Blog",
    id: `${baseUrl}/blog`,
    link: `${baseUrl}/blog`,
    language: "en",

    image: `${baseUrl}/banner.png`,
    favicon: `${baseUrl}/icon.png`,
    copyright: "Barnaby Keene",
  });

  const pages = blog.getPages();

  const sorted = pages.sort((a, b) => {
    return new Date(b.data.date!).getTime() - new Date(a.data.date!).getTime();
  });

  for (const page of sorted) {
    feed.addItem({
      id: page.url,
      title: page.data.title,
      description: page.data.description,
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
