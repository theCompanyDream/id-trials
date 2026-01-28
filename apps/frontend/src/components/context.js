import { create } from 'zustand'

const useUserStore = create((set, get) => ({
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20,
    userId: "uuid4",
    idTypesMap: {
        uuid4: { name: "UUID4", value: "uuid4", table: "users_uuid", analytics: "UUID" },
        cuidId: { name: "CUID", value: "cuidId", table: "users_cuid", analytics: "CUID" },
        snowId: { name: "SNOW", value: "snowId", table: "users_snowflake", analytics: "Snowflake" },
        ksuidId: { name: "KSUID", value: "ksuidId", table: "users_ksuid", analytics: "KSUID" },
        ulidId: { name: "ULID", value: "ulidId", table: "users_ulid", analytics: "ULID" },
        nanoId: { name: "NANO ID", value: "nanoId", table: "users_nanoid", analytics: "NanoID" }
    },
    getIdTypesArray: (renderfunc) => {
        const { idTypesMap } = get();
        return Object.entries(idTypesMap).map(([key, value]) => renderfunc(key, value));
    },
    updateStore: (filters) => set({ ...filters }),
    updateIdTypes: (newIdTypes) => {
        const { idTypesMap } = get()
        if (!idTypesMap[newIdTypes]) return;
        set({ userId: idTypesMap[newIdTypes].value })
    }
}))

export default useUserStore;