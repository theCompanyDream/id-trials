import { useEffect, useState } from "react";

import { useUserStore, Loading, PieChartComponent } from '@backpack';

const Analysis = () => {
	const [tableSize, setTableSize] = useState()

	useEffect(() => {
		fetch(`/analytics/tableSize?`)
			.then(data => data.json())
			.then(table => setTableSize(table))

	}, []);

	if (!tableSize) {
		return (
		<Loading />
		)
	}

	return (
		<main>
			<h1>Analysis Page</h1>

			{tableSize && (
				<PieChartComponent
					data={Object.entries(tableSize).map(([name, value]) => ({ name, value }))}
				/>)}
		</main>
	);
}

export default Analysis;