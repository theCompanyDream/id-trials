import React, { useState, memo } from 'react';
import { Link } from "react-router-dom";

interface LinkProp {
  link: string;
  name: string;
  img?: ImageBitmap;
}

interface NavigationProps {
  links: LinkProp[];
}

const Navigation = memo(({links}: NavigationProps) => {
  const [isOpen, setIsOpen] = useState(false);

  const isImageDesktop = (link: LinkProp) => {
    if (link.img) {
      return (
        <Link to={link.link} className="hover:text-yellow-300">
          <img src={link.img} alt="github logo" className="w-7 h-7" />
        </Link>
      );
    } else {
      return (
        <Link to={link.link} className="hover:text-yellow-300">
          {link.name}
        </Link>
      );
    }
  }

  const isImageMobile = (link: LinkProp) => {
    if (link.img) {
      return (
        <Link to={link.link} className="block px-3 py-2 rounded-md text-base font-medium hover:bg-yellow-600">
          <img src={link.img} alt="github logo" className="w-7 h-7" />
        </Link>
      );
    } else {
      return (
        <Link to={link.link} className="block px-3 py-2 rounded-md text-base font-medium hover:bg-yellow-600">
          {link.name}
        </Link>
      );
    }
  }

  return (
    <header className="bg-rose-500 text-white shadow w-full">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link to="/" className="text-xl font-bold">Id Trials</Link>

        <nav className="hidden md:flex space-x-4">
          {links && (links.map((link: LinkProp) => isImageDesktop(link)))}
        </nav>

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
      {isOpen && links && (links.map((link: LinkProp) => isImageMobile(link)))}
    </header>
  );
});

export default Navigation;
