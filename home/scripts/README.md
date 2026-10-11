# Homepage statistics

The homepage refreshes GitHub stars, commits and contributors daily through Next.js fetch revalidation. Commits and contributors use GitHub's last-page link with one result per page, so the count includes every page. Each request falls back independently to `src/lib/home-stats.snapshot.json` if GitHub is unavailable or rate-limits the deployment.

`pnpm build` counts HTTP operations in `api/openapi.yaml` before building the site. When [tokei](https://github.com/XAMPPRocky/tokei) is installed, it also measures lines of code using `tokei --output json --exclude home/src/lib/home-stats.snapshot.json .` from the repository root. This follows the original table's `Total.code` measurement and excludes blanks and comments. It includes generated code and respects repository ignore files.

Build hosts without tokei retain the last measured line count, with its measurement date recorded in the snapshot. They do not need Git history or a GitHub token to render all five figures.

To refresh the checked-in measurements, install tokei, run `pnpm stats:refresh` from `home`, and include the updated snapshot in the change. This also refreshes the GitHub fallback values; it fails if any GitHub request fails rather than dating old values as new measurements.

Run the loader's pagination and failure-handling tests with `pnpm exec tsx --test src/lib/home-stats.test.ts`.
