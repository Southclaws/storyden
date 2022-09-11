import { useForm } from "react-hook-form";

export function useAuthScreen() {
  const { register, handleSubmit } = useForm({});

  const onSubmit = async () => {
    const response = await fetch(
      `https://1f91-212-82-91-46.ngrok.io/api/v1/auth/webauthn/make/southclaws`,
      {
        method: "POST",
      }
    );

    if (!response.ok) {
      throw new Error("Request failed" + response.statusText);
    }

    const { publicKey } = await response.json();

    if (!publicKey) {
      throw new Error("Response was empty");
    }

    console.log("WEBAUTHN", publicKey);

    const publicKeyCredentialCreationOptions = {
      ...publicKey,

      rp: {
        id: "1f91-212-82-91-46.ngrok.ios",
        name: "Storyden",
      },

      // overwrite challenge and user with the correct format
      challenge: Uint8Array.from(publicKey.challenge as string, (c) =>
        c.charCodeAt(0)
      ),
      user: {
        id: Uint8Array.from(publicKey.user.id as string, (c) =>
          c.charCodeAt(0)
        ),
        name: publicKey.user.name,
        displayName: publicKey.user.displayName,
      },
    };

    const credential = await navigator.credentials.create({
      publicKey: publicKeyCredentialCreationOptions,
    });

    const creds = await fetch(
      `https://1f91-212-82-91-46.ngrok.io/api/v1/auth/webauthn/make`,
      {
        method: "GET",
      }
    );

    console.log(credential, creds);
  };

  return {
    register,
    onSubmit: handleSubmit(onSubmit),
  };
}
