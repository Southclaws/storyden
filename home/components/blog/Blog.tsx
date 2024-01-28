import Link from "next/link";
import { Page } from "nextra";
import { getPagesUnderRoute } from "nextra/context";
import styles from "./blog.module.css";
import { formatDistanceToNow } from "date-fns";

type Post = Page & {
  frontMatter: {
    description?: string;
    date?: Date;
  };
};

export function Blog() {
  const posts = getPagesUnderRoute("/blog").filter(
    (p) => p.name !== "index"
  ) as Post[];

  return (
    <ul className={styles["list"]}>
      {posts.map((post) => {
        const {
          route,
          meta: { title },
          frontMatter: { description, date },
        } = post;

        return (
          <li key={route}>
            <article>
              <Link href={route}>
                <h1>{title}</h1>
              </Link>
              <p className={styles["timestamp"]}>
                Posted{" "}
                <time>{formatDistanceToNow(date, { addSuffix: true })}</time>
              </p>
              <p>{description}</p>
            </article>
          </li>
        );
      })}
    </ul>
  );
}
