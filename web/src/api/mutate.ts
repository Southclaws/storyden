import type { ScopedMutator } from "swr";

type Target<K, T, R> = {
  key: K;
  optimistic: (current: T | undefined) => T | undefined;
  commit?: (current: T | undefined, result: R) => T | undefined;
};

type TxOptions = {
  revalidate?: boolean;
};

/**
 * Perform a transactional mutation with optimistic updates. If the action
 * fails, all changes are rolled back to their previous state. If the action
 * succeeds, the changes are committed to SWR's cache based on Target list.
 */
export async function mutateTransaction<
  R,
  const _Targets extends readonly Target<any, any, R>[],
>(
  mutate: ScopedMutator,
  targets: _Targets,
  action: () => Promise<R>,
  options: TxOptions = { revalidate: false },
): Promise<R> {
  // snapshot current values for rollback
  const snapshots = await Promise.all(
    targets.map(async (t) => ({
      target: t,
      // mutate(key) returns current cached data (Promise<T | undefined>)
      prev: await mutate(t.key),
    })),
  );

  // apply optimistic updates (no revalidate)
  await Promise.all(
    snapshots.map(({ target }) =>
      mutate(target.key, (current: any) => target.optimistic(current), {
        revalidate: false,
      }),
    ),
  );

  try {
    const result = await action();

    // commit using server result (or derived result)
    await Promise.all(
      snapshots.map(({ target }) =>
        mutate(
          target.key,
          (current: any) =>
            target.commit ? target.commit(current, result) : current,
          { revalidate: false },
        ),
      ),
    );

    if (options.revalidate) {
      await Promise.all(snapshots.map(({ target }) => mutate(target.key)));
    }

    return result;
  } catch (err) {
    // rollback
    await Promise.all(
      snapshots.map(({ target, prev }) =>
        mutate(target.key, prev, { revalidate: false }),
      ),
    );
    throw err;
  }
}
