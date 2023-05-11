import { WEB_ADDRESS } from "src/config";

export function getPermalinkForThread(slug: string) {
  return `${WEB_ADDRESS}/t/${slug}`;
}

export function getPermalinkForPost(slug: string, post: string) {
  return `${WEB_ADDRESS}/t/${slug}#${post}`;
}
