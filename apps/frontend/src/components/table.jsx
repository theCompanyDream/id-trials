import React, { useState, useEffect } from 'react';
import { Link } from "react-router-dom";


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
              Username
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              First Name
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Last Name
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Email
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider"
            >
              Department
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {users && users.length > 0 ? (
            users.map((user, index) => (
              <tr key={index}>
                <td className="px-6 py-4 whitespace-nowrap">{user.user_name}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.first_name}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.last_name}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.email}</td>
                <td className="px-6 py-4 whitespace-nowrap">{user.department || 'N/A'}</td>
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

const Pagination = ({ currentPage, totalPages, onPageChange }) => {
  const [pages, setPages] = useState([]);

  useEffect(() => {
    let newPages = [];

    if (totalPages <= 10) {
      // If there are 10 or fewer pages, show them all.
      newPages = Array.from({ length: totalPages }, (_, i) => i + 1);
    } else {
      if (currentPage <= 6) {
        // Show first 10 pages.
        newPages = Array.from({ length: 10 }, (_, i) => i + 1);
      } else if (currentPage + 4 >= totalPages) {
        // Show last 10 pages.
        newPages = Array.from({ length: 10 }, (_, i) => totalPages - 9 + i);
      } else {
        // Center the current page in the middle of the pagination.
        newPages = Array.from({ length: 10 }, (_, i) => currentPage - 5 + i);
      }
    }

    setPages(newPages);
  }, [currentPage, totalPages]);


  return (
    <div className="flex justify-center items-center mt-4 space-x-2">
      {/* Previous Button */}
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="px-4 py-2 border rounded disabled:opacity-50"
      >
        Prev
      </button>
      {/* Page Numbers */}
      {pages.map((page) => (
        <button
          key={page}
          onClick={() => onPageChange(page)}
          className={`px-4 py-2 border rounded ${
            page === currentPage
              ? 'bg-blue-500 text-white'
              : 'bg-white text-blue-500 hover:bg-blue-100'
          }`}
        >
          {page}
        </button>
      ))}
      {/* Next Button */}
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="px-4 py-2 border rounded disabled:opacity-50"
      >
        Next
      </button>
    </div>
  );
};

export default Table;
