import React, { useContext, useState, useMemo } from 'react';
import { UserContext, Table } from '../components';

const UserTable = () => {
  const { users, setUsers } = useContext(UserContext);
  const [isfetch, setFetched] = useState(false);
  const [search, setSearch] = useState("");

  // Function to fetch users with search and page parameters
  const fetchUsers = (page = 1, query = search) => {
    fetch(`/api/uuid4?search=${encodeURIComponent(query)}&page=${page}`)
      .then((response) => response.json())
      .then((data) => {
        setUsers(data);
        setFetched(true);
      })
      .catch((err) => console.error("Error fetching users:", err));
  };

  const onDelete = (userId) => {
    fetch(`/api/uuid4/${userId}`, {
      method: "DELETE"
    })
    .then((data) => {
      const newUsers = users.users.filter(user => userId != user.id);
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

  // Trigger initial data fetch if no users yet
  useMemo(() => {
    if (!isfetch) {
      fetchUsers();
      setFetched(true);
    }
  }, [isfetch, fetchUsers, setFetched]);

  return (
    <main>
      <header className="flex justify-between items-center p-4 bg-gray-100">
        <h2 className="text-2xl font-bold">User List</h2>
        <div className="flex items-center">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search users..."
            className="border border-gray-300 rounded-md px-3 py-1 focus:outline-none focus:border-blue-500"
          />
          <button
            onClick={handleSearch}
            className="ml-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md px-4 py-2"
          >
            Search
          </button>
        </div>
      </header>
      {users && (
        <Table
          users={users.users}
          currentPage={users.page}
          totalPages={users.page_count}
          onPageChange={onPageChange}
          onDelete={onDelete}
        />
      )}
    </main>
  );
};

export default UserTable;
