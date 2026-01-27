import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { Loading, PieChartComponent, StackedBarChart, TimeSeriesChart, useUserStore } from '@backpack';

const Analysis = () => {
	const idTypes = useUserStore((state) => state.idTypes);
	const updateIdTypes = useUserStore((state) => state.updateIdTypes);
	const userId = useUserStore((state) => state.userId);

	const [slider, setSlider] = useState(24);
	const [tableSize, setTableSize] = useState<any>(null);
	const [idEfficiency, setIdEfficiency] = useState<any>(null);
	const [comparison, setComparison] = useState<any>(null);
	const [percentile, setPercentile] = useState<any>(null);
	const [timeSeries, setTimeSeries] = useState<any>(null);

	// ← BUG FIX: This was overwriting itself!
	const fetchIdAnalytics = async (id: string, hour: number) => {
		try {
			// Fetch all in parallel
			const [detailsRes, timeSeriesRes, percentilesRes] = await Promise.all([
				fetch(`/api/analytics/${id}/details`),
				fetch(`/api/analytics/${id}/timeseries?hours=${hour}`),
				fetch(`/api/analytics/${id}/percentiles`)
			]);

			const details = await detailsRes.json();
			const timeSeriesData = await timeSeriesRes.json();
			const percentilesData = await percentilesRes.json();

			// Set state properly (was overwriting before!)
			setTimeSeries(timeSeriesData);
			setPercentile(percentilesData);

			updateIdTypes(id);
		} catch (error) {
			console.error('Failed to fetch analytics:', error);
		}
	};

	const onChangeSlider = (e: React.ChangeEvent<HTMLInputElement>) => {
		const hour = parseInt(e.target.value);
		setSlider(hour);
		if (userId) {
			fetchIdAnalytics(userId, hour);
		}
	};

	useEffect(() => {
		const fetchInitialData = async () => {
			try {
				await fetch(`/api/analytics/tableSize?`)
					.then(data => data.json())
					.then(table => setTableSize(table.map(t => ({ name: t.table_name, value: t.size, sizePretty: t.size_pretty }))))

				await fetch(`/api/analytics/idEfficiency`)
					.then(data => data.json())
					.then(efficiency => setIdEfficiency(efficiency))

				await fetch("/api/analytics/comparison")
					.then(data => data.json())
					.then(comparison => setComparison(comparison))
			} catch (error) {
				console.error('Failed to fetch initial data:', error);
			}
		};

		fetchInitialData();
	}, []);

	if (!tableSize) {
		return <Loading />;
	}

	return (
		<main className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
			{/* Header */}
			<section className="bg-white shadow-sm border-b">
				<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
					<h1 className="text-3xl font-bold text-gray-900">ID Performance Analytics</h1>
					<p className="mt-2 text-gray-600">
						Comprehensive analysis of different ID generation strategies
					</p>
				</div>
			</section>

			<section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">

				{/* Overview Cards */}
				<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
					<div className="bg-white rounded-lg shadow p-6 border-l-4 border-blue-500">
						<div className="flex items-center">
							<div className="flex-shrink-0">
								<svg className="h-8 w-8 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
								</svg>
							</div>
							<div className="ml-4">
								<p className="text-sm font-medium text-gray-600">Total Tables</p>
								<p className="text-2xl font-bold text-gray-900">{tableSize?.length || 0}</p>
							</div>
						</div>
					</div>

					<div className="bg-white rounded-lg shadow p-6 border-l-4 border-green-500">
						<div className="flex items-center">
							<div className="flex-shrink-0">
								<svg className="h-8 w-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
								</svg>
							</div>
							<div className="ml-4">
								<p className="text-sm font-medium text-gray-600">Most Efficient</p>
								<p className="text-2xl font-bold text-gray-900">
									{idEfficiency?.[0]?.table_name?.replace('users_', '').toUpperCase() || 'N/A'}
								</p>
							</div>
						</div>
					</div>

					<div className="bg-white rounded-lg shadow p-6 border-l-4 border-purple-500">
						<div className="flex items-center">
							<div className="flex-shrink-0">
								<svg className="h-8 w-8 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
								</svg>
							</div>
							<div className="ml-4">
								<p className="text-sm font-medium text-gray-600">Time Range</p>
								<p className="text-2xl font-bold text-gray-900">{slider}h</p>
							</div>
						</div>
					</div>
				</div>

				{/* Main Charts Grid */}
				<div className="grid grid-cols-1 lg:grid-cols-2 gap-8">

					{/* Table Size Chart */}
					{tableSize && idEfficiency && (
						<div className="bg-white rounded-lg shadow-lg p-6">
							<div className="flex items-center mb-4">
								<div className="h-2 w-2 rounded-full bg-blue-500 mr-2"></div>
								<h2 className="text-xl font-bold text-gray-900">Storage Distribution</h2>
							</div>
							<p className="text-sm text-gray-600 mb-6">
								Database storage allocation across ID types
							</p>
							<PieChartComponent
								data={tableSize.map((t: any) => ({
									name: t.name,
									value: t.value,
									valuePretty: t.sizePretty
								}))}
								data2={idEfficiency.map((d: any) => ({
									name: d.table_name,
									value: d.row_count,
									valuePretty: `${d.row_count.toLocaleString()} rows`
								}))}
								labelLine={true}
								width="100%"
								height={450}
							/>
						</div>
					)}

					{/* Efficiency Comparison */}
					{idEfficiency && (
						<div className="bg-white rounded-lg shadow-lg p-6">
							<div className="flex items-center mb-4">
								<div className="h-2 w-2 rounded-full bg-green-500 mr-2"></div>
								<h2 className="text-xl font-bold text-gray-900">Efficiency Metrics</h2>
							</div>
							<p className="text-sm text-gray-600 mb-6">
								Storage efficiency and waste factor analysis
							</p>
							<StackedBarChart
								data={idEfficiency.map(d => ({
									name: d.table_name.replace('users_', '').toUpperCase(),
									theoretical: d.theoretical_min_bytes,
									wasted: d.avg_id_bytes - d.theoretical_min_bytes,
								}))}
								xAxisKey="name"
								stacks={[
									{ dataKey: 'theoretical', name: 'Minimum Required', fill: '#192fbc' },
									{ dataKey: 'wasted', name: 'Wasted Space', fill: '#ef4444' }
								]}
								yAxisLabel="Bytes per ID"
								showGrid={true}
								height={450}
							/>
						</div>
					)}
				</div>

				{/* {idEfficiency && (
					<div className="bg-white rounded-lg shadow-lg p-6">
						<div className="flex items-center mb-4">
							<div className="h-2 w-2 rounded-full bg-green-500 mr-2"></div>
							<h2 className="text-xl font-bold text-gray-900">Efficiency Metrics</h2>
						</div>
						<p className="text-sm text-gray-600 mb-6">
							Storage efficiency and waste factor analysis
						</p>
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
						/>
					</div>)} */}

				{/* Performance Comparison - Full Width */}
				{comparison && (
					<div className="bg-white rounded-lg shadow-lg p-6">
						<div className="flex items-center mb-4">
							<div className="h-2 w-2 rounded-full bg-purple-500 mr-2"></div>
							<h2 className="text-xl font-bold text-gray-900">Performance Comparison</h2>
						</div>
						<p className="text-sm text-gray-600 mb-6">
							Average response time and request volume by ID type
						</p>
						<TimeSeriesChart
							data={comparison}
							series={[
								{ dataKey: 'avg_duration', name: 'Avg Duration (ms)', stroke: '#10b981' },
								{ dataKey: 'request_count', name: 'Request Count', stroke: '#3b82f6' },
							]}
							xAxisKey="id_type"
							height={450}
						/>
					</div>
				)}

				{/* Time Range Slider */}
				<div className="bg-white rounded-lg shadow-lg p-6">
					<div className="flex items-center justify-between mb-6">
						<div>
							<h3 className="text-lg font-semibold text-gray-900">Time Range Filter</h3>
							<p className="text-sm text-gray-600">
								Adjust the time window for detailed analytics
							</p>
						</div>
						<div className="bg-blue-100 text-blue-800 px-4 py-2 rounded-lg font-bold">
							{slider} hours
						</div>
					</div>
					<input
						type="range"
						min="1"
						max="24"
						value={slider}
						onChange={onChangeSlider}
						className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-blue-500"
					/>
					<div className="flex justify-between text-xs text-gray-500 mt-2">
						<span>1h</span>
						<span>6h</span>
						<span>12h</span>
						<span>18h</span>
						<span>24h</span>
					</div>
				</div>

				{/* ID Type Selection */}
				<div className="bg-white rounded-lg shadow-lg p-6">
					<div className="flex items-center mb-6">
						<div className="h-2 w-2 rounded-full bg-indigo-500 mr-2"></div>
						<h3 className="text-lg font-semibold text-gray-900">Detailed ID Analysis</h3>
					</div>
					<p className="text-sm text-gray-600 mb-6">
						Select an ID type to view detailed performance metrics
					</p>
					<div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
						{idTypes.map((idType) => (
							<button
								key={idType.value}
								onClick={() => fetchIdAnalytics(idType.value, slider)}
								className="group relative overflow-hidden bg-gradient-to-br from-blue-500 to-blue-600 text-white p-4 rounded-lg shadow-md hover:shadow-xl hover:from-blue-600 hover:to-blue-700 transition-all duration-200 transform hover:scale-105"
							>
								<div className="flex flex-col items-center">
									<svg className="h-8 w-8 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
									</svg>
									<span className="font-semibold text-sm">{idType.name}</span>
								</div>
								<div className="absolute inset-0 bg-white opacity-0 group-hover:opacity-10 transition-opacity"></div>
							</button>
						))}
					</div>
				</div>

				{/* Time Series and Percentiles - Conditional Rendering */}
				{timeSeries && (
					<div className="bg-white rounded-lg shadow-lg p-6">
						<div className="flex items-center mb-4">
							<div className="h-2 w-2 rounded-full bg-orange-500 mr-2"></div>
							<h3 className="text-lg font-semibold text-gray-900">Time Series Data</h3>
						</div>
						<TimeSeriesChart
							data={timeSeries}
							series={[
								{ dataKey: 'avg_duration', name: 'Average Duration', stroke: '#f97316' },
							]}
							xAxisKey="timestamp"
							height={400}
						/>
					</div>
				)}

				{percentile && (
					<div className="bg-white rounded-lg shadow-lg p-6">
						<div className="flex items-center mb-4">
							<div className="h-2 w-2 rounded-full bg-teal-500 mr-2"></div>
							<h3 className="text-lg font-semibold text-gray-900">Percentile Distribution</h3>
						</div>
						<div className="grid grid-cols-2 md:grid-cols-5 gap-4">
							{['p50', 'p75', 'p90', 'p95', 'p99'].map((p) => (
								<div key={p} className="bg-gray-50 rounded-lg p-4 text-center">
									<p className="text-sm text-gray-600 uppercase mb-1">{p}</p>
									<p className="text-2xl font-bold text-gray-900">
										{percentile[p]?.toFixed(2) || 'N/A'}ms
									</p>
								</div>
							))}
						</div>
					</div>
				)}

				{/* Call to Action */}
				<div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg shadow-lg p-8 text-center">
					<h3 className="text-2xl font-bold text-white mb-4">
						Want to explore the raw data?
					</h3>
					<p className="text-blue-100 mb-6">
						Dive deeper into the database with our interactive data explorer
					</p>
					<Link
						to="/explore"
						className="inline-flex items-center px-6 py-3 bg-white text-blue-600 font-semibold rounded-lg shadow-md hover:bg-gray-50 transition-colors duration-200"
					>
						<svg className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
						</svg>
						Explore Data
					</Link>
				</div>

			</section>
		</main>
	);
}

export default Analysis;