import { create } from 'zustand'

const useUserStore = create((set, get) => ({
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20,
    userId: "uuid",
    updateStore: (filters) => set({ ...filters }),
}))

export default useUserStore;