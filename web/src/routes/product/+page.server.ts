import type { PageServerLoad } from './$types'

export const load: PageServerLoad<{ a: number }> = () => {
  return { a: 123 }
}
