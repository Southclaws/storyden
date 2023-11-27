import { LinkWithRefs } from "src/api/openapi/schemas";
import { Link } from "src/theme/components/Link";

import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  link: LinkWithRefs;
};

export function LinkView({ link }: Props) {
  return (
    <Flex flexDir="column">
      <styled.h1 fontSize="lg">{link.title}</styled.h1>
      <p>{link.description}</p>
      <a href={link.url}>{link.url}</a>
      <Link href={`/l?q=${link.domain}`}>{link.domain}</Link>
      <pre>{link.slug}</pre>
      {link.threads.map((v) => (
        <>{v.title}</>
      ))}
      {link.assets.map((v) => (
        <>{v.url}</>
      ))}
    </Flex>
  );
}
