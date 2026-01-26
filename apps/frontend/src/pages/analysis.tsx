import { useEffect, useState } from "react";

import { useUserStore, Loading, PieChartComponent } from '@backpack';

const Analysis = () => {
	const [tableSize, setTableSize] = useState()

	useEffect(() => {
		fetch(`/api/analytics/tableSize?`)
			.then(data => data.json())
			.then(table => setTableSize(table.map(t => ({ name: t.table_name, value: t.size, sizePretty: t.size_pretty }))))

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
					data={tableSize}
					title="Table Size Analysis"
					legendPosition="bottom"
					width="30%"
				/>)}
		</main>
	);
}

export default Analysis;