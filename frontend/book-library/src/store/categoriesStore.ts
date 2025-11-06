import { create } from 'zustand'
import { getCategoriesWithBooks, getUncategorizedBooks } from '../api'

interface CategoryStore {
	categories: any[]
	uncategorizedBooks: any[]
	loading: boolean

	refresh: () => Promise<void>
}

export const useCategoriesStore = create<CategoryStore>((set: any) => ({
	categories: [],
	uncategorizedBooks: [],
	loading: false,

	refresh: async () => {
		set({ loading: true })

		const [cats, uncats] = await Promise.all([
			getCategoriesWithBooks(),
			getUncategorizedBooks(),
		])

		set({
			categories: Array.isArray(cats) ? cats : [],
			uncategorizedBooks: Array.isArray(uncats) ? uncats : [],
			loading: false,
		})
	},
}))
