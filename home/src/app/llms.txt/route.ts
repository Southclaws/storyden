import { source } from '@/lib/source';
import { llms } from 'fumadocs-core/source';

// cached forever
export const revalidate = false;

export function GET() {
  return new Response(llms(source).index());
}
