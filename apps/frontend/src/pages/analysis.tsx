import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { Loading, PieChartComponent, StackedBarChart } from '@backpack';

const Analysis = () => {
	const [tableSize, setTableSize] = useState()
	const [idEfficiency, setIdEfficiency] = useState()
	const [percentile, setPercentile] = useState<string | null>(null)

	const fetchIdPercentile = (id: string) => {
		fetch(`/analytics/${id}/percentiles`)
			.then(data => data.json())
			.then(efficiency => setPercentile(efficiency))
	}

	useEffect(() => {
		fetch(`/api/analytics/tableSize?`)
			.then(data => data.json())
			.then(table => setTableSize(table.map(t => ({ name: t.table_name, value: t.size, sizePretty: t.size_pretty }))))

		fetch(`/api/analytics/idEfficiency`)
			.then(data => data.json())
			.then(efficiency => setIdEfficiency(efficiency))
	}, []);

	if (!tableSize) {
		return (
		<Loading />
		)
	}

	return (
		<main>
			{tableSize && (
				<PieChartComponent
					data={tableSize}
					title="Table Size Analysis"
					legendPosition="bottom"
					width="50%"
					height={450}
				/>)}

			{ idEfficiency && (
				<StackedBarChart
					data={idEfficiency}
					title="ID Efficiency Analysis"
					stacks={[
						{ dataKey: 'row_count', name: 'Row Count', fill: '#10b981' },
					]}
					xAxisLabel="table_name"
					showGrid={true}
					width="50%"
					height={450}
				/>)}

			{ idEfficiency && (
				<StackedBarChart
					data={idEfficiency}
					title="ID Efficiency Analysis"
					stacks={[
						{ dataKey: 'efficiency_percent', name: 'Efficiency %', fill: '#3b82f6' },
						{ dataKey: 'waste_factor', name: 'Wasted Space', fill: '#ef4444' },
						{ dataKey: 'theoretical_min_bytes', name: 'Minimum Required', fill: '#10b981' },
					]}
					xAxisLabel="table_name"
					showGrid={true}
					width="50%"
					height={450}
				/>)}

			<section className="flex justify-center mt-8">
				<Link to="/explore" className="text-blue-500 text-center bg-blue-500 p-3 text-white rounded-md hover:bg-blue-600">
					Explore More Data
				</Link>
			</section>
		</main>
	);
}

export default Analysis;