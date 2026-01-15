import React, { useState, useEffect } from 'react';
import { useUserStore, Table } from '../components';

const UserTable = () => {
  const users = useUserStore((state) => state.users)
  const page = useUserStore((state) => state.page)
  const page_count = useUserStore((state) => state.page_count)
  const userId = useUserStore((state) => state.userId)
  const updateStore = useUserStore((state) => state.updateStore)
  const [isfetch, setFetched] = useState(false);
  const [search, setSearch] = useState("");

  // Function to fetch users with search and page parameters
  const fetchUsers = (page = 1, query = search) => {
    fetch(`/api/${userId}?search=${encodeURIComponent(query)}&page=${page}`)
      .then((response) => response.json())
      .then((data) => {
        updateStore({...data})
        setFetched(true);
      })
      .catch((err) => console.error("Error fetching users:", err));
  };

  const onDelete = (id) => {
    fetch(`/api/${userId}/${id}`, {
      method: "DELETE"
    })
    .then((data) => {
      const newUsers = users.users.filter(user => id != user.id);
      setUsers({...users, users: newUsers})
    })
  }

  // Handler for page changes
  const onPageChange = (page) => {
    fetchUsers(page);
  };

  // Handler for search button click
  const handleSearch = () => {
    // Start a fresh search on page 1
    fetchUsers(1, search);
  };

  const handleSelect = (e) => {
    updateStore({
      userId: e.target.value
    })
  }

  // Trigger initial data fetch if no users yet
  useEffect(() => {
    if (!isfetch) {
      fetchUsers();
      setFetched(true);
    }
  }, [isfetch, fetchUsers, setFetched]);

  return (
    <main>
      <header className="flex justify-between items-center p-6">
        <h2 className="text-3xl font-bold text-white">User Directory</h2>
        <div className="flex items-center gap-3">
          <select
            onChange={handleSelect}  // Changed from onSelect to onChange
            value={userId}
            className="px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-800 font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition hover:border-blue-400 cursor-pointer"
          >
            <option value="uuid4">UUID</option>
            <option value="cuid">CUID</option>
            <option value="snow">Snowflake</option>
            <option value="ksuid">KSUID</option>
            <option value="ulid">ULID</option>
            <option value="nano">NanoID</option>
          </select>

          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search users..."
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition bg-white"
          />

          <button
            onClick={handleSearch}
            className="px-6 py-2 bg-green-500 hover:bg-green-600 active:bg-green-700 text-white font-medium rounded-lg transition shadow-md hover:shadow-lg"
          >
            Search
          </button>
        </div>
      </header>
      {users && (
        <Table
          users={users}
          currentPage={page}
          totalPages={page_count}
          onPageChange={onPageChange}
          onDelete={onDelete}
        />
      )}
    </main>
  );
};

export default UserTable;
