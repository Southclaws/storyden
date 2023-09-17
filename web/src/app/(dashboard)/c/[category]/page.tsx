import { FeedScreen } from "src/screens/feed/FeedScreen";

type Props = {
  params: {
    category: string;
  };
};

export default function Page(props: Props) {
  return <FeedScreen category={props.params.category} />;
}
