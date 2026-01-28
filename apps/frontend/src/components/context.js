import { create } from 'zustand'

const useUserStore = create((set, get) => ({
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20,
    userId: "uuid4",
    idTypes: [
        {name: "UUID4", value: "uuid4", table: "users_uuid", analytics: "UUID"},
        {name: "CUID", value: "cuidId", table: "users_cuid", analytics: "CUID"},
        {name: "SNOW", value: "snowId", table: "users_snowflake", analytics: "Snowflake"},
        {name: "KSUID", value: "ksuidId", table: "users_ksuid", analytics: "KSUID"},
        {name: "ULID", value: "ulidId", table: "users_ulid", analytics: "ULID"},
        {name: "NANO ID", value: "nanoId", table: "users_nanoid", analytics: "NanoID"}
    ],
    updateStore: (filters) => set({ ...filters }),
    updateIdTypes: (newIdTypes) => {
        if (!Array.isArray(newIdTypes)) return;
        set({ idTypes: newIdTypes })
    }
}))

export default useUserStore;