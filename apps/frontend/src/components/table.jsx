import { Link } from "react-router-dom";

import { Pagination } from "@backpack"


// Table component with pagination
const Table = ({ users, currentPage, totalPages, onPageChange, onDelete }) => (
  <section>
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Id
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Username
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Name
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Email
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {users && users.length > 0 ? (
            users.map((user, index) => (
              <tr key={index}>
                <td className="px-6 py-4 whitespace-nowrap">{String(user.id)}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.user_name}</td>
                <td className="px-6 py-4 whitespace-nowrap">{[user.first_name, user.last_name].join(" ")}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.email}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <Link to={`/detail/${user.id}`} className='bg-blue-500 text-white px-4 py-2 border rounded'>Edit</Link>
                  <button onClick={() => onDelete(user.id)} className='bg-red-500 text-white px-4 py-2 border rounded'>Delete</button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="6" className="px-6 py-4 text-center">
                No users found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
    {/* Pagination Controls */}
    <Pagination
      currentPage={currentPage}
      totalPages={totalPages}
      onPageChange={onPageChange}
    />
  </section>
);

export default Table;
