import { Post } from "src/api/openapi/schemas";

export function Post(props: Post) {
  return <>{props.body}</>;
}
