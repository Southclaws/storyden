import { useEffect, useState } from "react";

import { linkCreate } from "src/api/openapi/links";
import { Link } from "src/api/openapi/schemas";

export type Props = {
  href: string;
};

async function hydrateLink(url: string) {
  return await linkCreate({ url });
}

export function useRichLink({ href }: Props) {
  const [link, setLink] = useState<Link | undefined>(undefined);

  useEffect(() => {
    hydrateLink(href).then(setLink).catch();
  }, [href]);

  return {
    link,
  };
}
