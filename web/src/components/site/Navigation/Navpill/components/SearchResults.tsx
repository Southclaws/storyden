import { PostProps } from "src/api/openapi/schemas";

import { styled } from "@/styled-system/jsx";

type Props = {
  results: PostProps[];
};
export function SearchResults(props: Props) {
  return (
    <styled.ol m="0">
      {props.results.map((v) => (
        <styled.li key={v.id}>
          <p>{v.body}</p>
        </styled.li>
      ))}
    </styled.ol>
  );
}
