import { notFound, redirect } from "next/navigation";

import { DatagraphItemKind } from "@/api/openapi-schema";
import { WEB_ADDRESS } from "@/config";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export type Props = {
  params: Promise<{
    kind: string;
    mark: string;
  }>;
  searchParams: Promise<{
    [key: string]: string | string[] | undefined;
  }>;
};

function resolvePath(kind: string, mark: string): string | undefined {
  switch (kind) {
    case DatagraphItemKind.post:
      return `/t/locate/${mark}`;
    case DatagraphItemKind.thread:
      return `/t/${mark}`;
    case DatagraphItemKind.reply:
      return `/t/locate/${mark}`;
    case DatagraphItemKind.node:
      return `/l/${mark}`;
    case DatagraphItemKind.collection:
      return `/c/${mark}`;
    case DatagraphItemKind.profile:
      return `/m/${mark}`;
    default:
      return undefined;
  }
}

export default async function Page(props: Props) {
  const { kind, mark } = await props.params;
  const searchParams = await props.searchParams;

  const path = resolvePath(kind, mark);
  if (!path) {
    notFound();
  }

  const url = new URL(path, WEB_ADDRESS);

  Object.entries(searchParams).forEach(([key, value]) => {
    if (value === undefined) return;
    if (typeof value === "string") {
      url.searchParams.set(key, value);
    } else if (Array.isArray(value)) {
      value.forEach((v) => url.searchParams.append(key, v));
    }
  });

  redirect(url.toString(), "replace");
}
