import { WEB_ADDRESS } from "src/config";

export function getPermalinkForPost(
  slug: string,
  post: string,
  page?: number,
) {
  const url = new URL(`/t/${slug}`, WEB_ADDRESS);

  if (page && page > 1) {
    url.searchParams.set("page", page.toString());
  }

  url.hash = post;

  return url.toString();
}

export function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}
