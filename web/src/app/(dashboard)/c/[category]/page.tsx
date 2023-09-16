import { FeedScreen } from "src/screens/feed/FeedScreen";

type Props = { params: { c: string } };

export default function Page({ params: { c } }: Props) {
  return <FeedScreen category={c} />;
}
