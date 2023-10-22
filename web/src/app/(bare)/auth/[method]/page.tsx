"use client";

import { useParams } from "next/navigation";
import { z } from "zod";

import { AuthScreen } from "src/screens/auth/AuthScreen";

const ParamSchema = z.object({ method: z.string().optional() });

function Page() {
  const params = useParams();
  const { method } = ParamSchema.parse(params);

  return <AuthScreen method={method} />;
}

export default Page;
