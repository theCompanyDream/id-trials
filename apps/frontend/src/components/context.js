import { createContext } from 'react'

const UserObject = {
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20
}

export default createContext(UserObject)