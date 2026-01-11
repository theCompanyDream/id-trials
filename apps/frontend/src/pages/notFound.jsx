import React, {memo} from 'react';
import { Link } from 'react-router-dom'; // Optional: if you're using react-router

const NotFoundPage = memo(() => (
	<main className="min-h-screen flex flex-col items-center justify-center bg-gray-100 p-4">
		<h1 className="text-9xl font-extrabold">404</h1>
		<p className="mt-4 text-3xl font-bold">
		Oops! Page not found.
		</p>
		<p className="mt-2 text-center max-w-md">
		The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
		</p>
		<Link
		to="/"
		className="mt-6 inline-block px-6 py-3 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700 transition"
		>
		Go back Home
		</Link>
	</main>
));

export default NotFoundPage;
