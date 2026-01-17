import { useEffect, useState } from "react";

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: Function
}

const Pagination = ({ currentPage, totalPages, onPageChange }: PaginationProps) => {
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

export default Pagination;