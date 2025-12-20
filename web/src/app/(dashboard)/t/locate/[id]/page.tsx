import { redirect } from "next/navigation";

import { PostLocationKind } from "@/api/openapi-schema";
import { postLocationGet } from "@/api/openapi-server/posts";
import { WEB_ADDRESS } from "@/config";

export type Props = {
  params: Promise<{
    id: string;
  }>;
};

export default async function LocatePage(props: Props) {
  const { id } = await props.params;

  const { data } = await postLocationGet({ id });

  if (data.kind === PostLocationKind.thread) {
    redirect(`/t/${data.slug}`);
  }

  const url = new URL(`/t/${data.slug}`, WEB_ADDRESS);

  if (data.page && data.page > 1) {
    url.searchParams.set("page", data.page.toString());
  }

  url.hash = id;

  redirect(url.toString());
}
