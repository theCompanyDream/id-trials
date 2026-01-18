import { useEffect, useState } from "react";

const Analysis = () => {
	const [tableSize, setTableSize] = useState()

	useEffect(() => {
		fetch(`/analytics/tableSize?`)
			.then(data => data.json())
			.then(table => setTableSize(table))

	}, []);

	return (
		<main>Analytics</main>
	);
}

export default Analysis;