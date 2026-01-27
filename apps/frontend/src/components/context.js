import { create } from 'zustand'

const useUserStore = create((set, get) => ({
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20,
    userId: "uuid4",
    idTypes: [
        {name: "UUID4", value: "uuid4"},
        {name: "CUID", value: "cuidId"},
        {name: "SNOW", value: "snowId"},
        {name: "KSUID", value: "ksuidId"},
        {name: "ULID", value: "ulidId"},
        {name: "NANO ID", value: "nanoId"}
    ],
    updateStore: (filters) => set({ ...filters }),
}))

export default useUserStore;