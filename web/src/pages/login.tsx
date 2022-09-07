import { GetServerSideProps, NextPage } from "next";
import { AuthScreen } from "../screens/auth/AuthScreen";

type Props = {
  //
};

const Page: NextPage = (props: Props) => {
  return <AuthScreen {...props} />;
};

export const getServerSideProps: GetServerSideProps<Props> = async () => {
  return {
    props: {},
  };
};

export default Page;
