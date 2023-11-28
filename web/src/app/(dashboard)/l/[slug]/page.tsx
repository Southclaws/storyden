import { LinkScreen } from "src/screens/directory/links/LinkScreen/LinkScreen";

type Props = {
  params: {
    slug: string;
  };
};

export default function Page(props: Props) {
  return <LinkScreen slug={props.params.slug} />;
}
