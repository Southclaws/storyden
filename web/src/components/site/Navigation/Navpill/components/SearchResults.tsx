import { PostProps } from "src/api/openapi/schemas";
import { ListItem, OrderedList, Text } from "src/theme/components";

type Props = {
  results: PostProps[];
};
export function SearchResults(props: Props) {
  return (
    <OrderedList m={0}>
      {props.results.map((v) => (
        <ListItem key={v.id}>
          <Text>{v.body}</Text>
        </ListItem>
      ))}
    </OrderedList>
  );
}
