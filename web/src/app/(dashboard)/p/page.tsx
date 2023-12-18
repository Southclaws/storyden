import { MemberIndexScreen } from "src/screens/directory/members/MemberIndexScreen/MemberIndexScreen";

type Props = {
  searchParams: {
    q: string;
    page: number;
  };
};

export default function Page(props: Props) {
  return (
    <MemberIndexScreen
      query={props.searchParams.q}
      page={props.searchParams.page}
    />
  );
}
