import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { Loading, PieChartComponent, StackedBarChart, TimeSeriesChart, useUserStore } from '@backpack';

const Analysis = () => {
	const idTypes = useUserStore((state) => state.idTypes)
	const [chosenIdType, setChosenIdType] = useState()
	const [slider, setSlider] = useState(24)
	const [tableSize, setTableSize] = useState()
	const [idEfficiency, setIdEfficiency] = useState()
	const [comparison, setComparison] = useState<{ [key: string]: any }>({})
	const [percentile, setPercentile] = useState<{ [key: string]: any }>({})
	const [timeSeries, setTimeSeries] = useState<{ [key: string]: any }>({})

	const fetchIdAnalytics = async(id: string, hour: number) => {
		fetch(`/analytics/${id}/details`)
			.then(data => data.json())
			.then(efficiency => setTimeSeries({id: efficiency}))

		fetch(`/analytics/${id}/timeseries?hours=${hour}`)
			.then(data => data.json())
			.then(efficiency => setTimeSeries({id: efficiency}))

		fetch(`/analytics/${id}/percentiles`)
			.then(data => data.json())
			.then(efficiency => setPercentile({id: efficiency}))

		setChosenIdType(id)
	}

	const onChangeSlider = (e: any) => {
		const hour = parseInt(e.target.value)
		setSlider(hour)
		fetchIdAnalytics(chosenIdType, hour)
	}

	useEffect(() => {
		fetch(`/api/analytics/tableSize?`)
			.then(data => data.json())
			.then(table => setTableSize(table.map(t => ({ name: t.table_name, value: t.size, sizePretty: t.size_pretty }))))

		fetch(`/api/analytics/idEfficiency`)
			.then(data => data.json())
			.then(efficiency => setIdEfficiency(efficiency))

		fetch("/api/analytics/comparison")
			.then(data => data.json())
			.then(comparison => setComparison(comparison))
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

			{idEfficiency && (
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

			{idEfficiency && (
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

			{comparison && <TimeSeriesChart
				data={comparison}
				title="ID Types Comparison Over Time"
				series={[
					{ dataKey: 'id_type', name: 'UUID', stroke: '#10b981' },
				]}
				xAxisLabel="timestamp"
				yAxisLabel="count"
				width="80%"
				height={450}
			/>}

			<input type="range" min="1" max="24" defaultValue="24" value={slider} onChange={onChangeSlider} />

			<ul>
				{idTypes.map((idType) => (
					<li key={idType.value} className="mb-8">
						<button onClick={() => fetchIdAnalytics(idType.value, slider)} className="text-blue-500 text-center bg-blue-500 p-3 text-white rounded-md hover:bg-blue-600">{idType.name} Analysis</button>
					</li>
				))}
			</ul>



			<section className="flex justify-center mt-8">
				<p>
					If you want to explore more data, click the button below:
					&nbsp;
				</p>
				<Link to="/explore" className="text-blue-500 text-center bg-blue-500 p-3 text-white rounded-md hover:bg-blue-600">
					Explore More Data
				</Link>
			</section>
		</main>
	);
}

export default Analysis;