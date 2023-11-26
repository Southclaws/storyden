import { useEffect, useState } from "react";

import { linkCreate } from "src/api/openapi/links";
import { Link } from "src/api/openapi/schemas";

export type Props = {
  href: string;
  initial?: Link;
};

async function hydrateLink(url: string) {
  return await linkCreate({ url });
}

export function useRichLink({ href, initial }: Props) {
  const [link, setLink] = useState<Link | undefined>(initial);

  useEffect(() => {
    hydrateLink(href).then(setLink).catch();
  }, [href]);

  return {
    link,
  };
}
