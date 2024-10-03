import { LinkScreen } from "src/screens/library/links/LinkScreen/LinkScreen";

type Props = {
  params: {
    slug: string;
  };
};

export default function Page(props: Props) {
  return <LinkScreen slug={props.params.slug} />;
}
