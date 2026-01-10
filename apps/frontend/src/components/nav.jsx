import React, { useState, memo } from 'react';
import { Link } from "react-router-dom";
import github from "../assets/github.svg"

const Navigation = memo(() => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="bg-gray-800 text-white shadow w-full">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        {/* Logo or Title */}
        <Link to="/" className="text-xl font-bold">User Administrator</Link>
        {/* Desktop Navigation */}
        <nav className="hidden md:flex space-x-4">
          <Link to="/detail" className="hover:text-gray-300">Create</Link>
          <Link to="/about" className="hover:text-gray-300">About</Link>
          <Link to="https://user.tbrantleyii.dev/api/swagger/index.html" className="hover:text-gray-300">Docs</Link>
          <Link to="https://github.com/theCompanyDream/id-trials" className="hover:text-gray-300">
            <img src={github} alt="github logo" class="w-7 -c-7" />
          </Link>
        </nav>
        {/* Mobile Menu Button */}
        <div className="md:hidden">
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="focus:outline-none"
          >
            {isOpen ? (
              // Close (X) icon
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            ) : (
              // Hamburger icon
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            )}
          </button>
        </div>
      </div>
      {/* Mobile Navigation */}
      {isOpen && (
        <div className="md:hidden bg-gray-700">
          <nav className="px-4 py-2 space-y-1">
            <Link to="/detail" className="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-600">
              Create
            </Link>
            <Link to="/services" className="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-600">
              Docs
            </Link>
            <Link to="/about" className="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-600">
              About
            </Link>
            <Link to="https://github.com/theCompanyDream/id-trials" className="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-600">
              <img src={github} alt="github logo" />
            </Link>
          </nav>
        </div>
      )}
    </header>
  );
});

export default Navigation;
