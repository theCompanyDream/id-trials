import React from 'recharts';

interface BoxPlotData {
  route_path: string;
  http_method: string;
  min_duration: number;
  quartile1: number;
  median: number;
  quartile3: number;
  max_duration: number;
  [key: string]: any; // Other fields
}

interface BoxPlotProps {
  data: BoxPlotData[];
  width?: number | string;
  height?: number | string;
  yAxisLabel?: string;
}

const BoxPlotChart: React.FC<BoxPlotProps> = ({
  data,
  width = '100%',
  height = 400,
  yAxisLabel = 'Duration (ms)'
}) => {
  // Transform data for rendering
  const transformedData = data.map((item, idx) => ({
    id: idx,
    label: `${item.route_path}\n${item.http_method}`,
    min: item.min_duration,
    q1: item.quartile1,
    median: item.median,
    q3: item.quartile3,
    max: item.max_duration
  }));

  const yMin = Math.min(...data.map(d => d.min_duration)) * 0.9;
  const yMax = Math.max(...data.map(d => d.max_duration)) * 1.1;
  const yRange = yMax - yMin;

  const boxWidth = 60;
  const padding = 80;

  return (
    <svg width={width} height={height} style={{ border: '1px solid #e5e7eb' }}>
      {/* Y-axis */}
      <line x1={padding} y1={20} x2={padding} y2={height - 40} stroke="#d1d5db" />

      {/* X-axis */}
      <line x1={padding} y1={height - 40} x2={width} y2={height - 40} stroke="#d1d5db" />

      {/* Y-axis label */}
      <text
        x={20}
        y={height / 2}
        textAnchor="middle"
        transform={`rotate(-90 20 ${height / 2})`}
        className="text-sm font-medium fill-gray-600"
      >
        {yAxisLabel}
      </text>

      {/* Grid lines & ticks */}
      {[0, 0.25, 0.5, 0.75, 1].map((tick, i) => {
        const value = yMin + tick * yRange;
        const y = height - 40 - (tick * (height - 60));
        return (
          <g key={i}>
            <line x1={padding - 5} y1={y} x2={padding} y2={y} stroke="#9ca3af" />
            <text x={padding - 10} y={y + 4} textAnchor="end" className="text-xs fill-gray-500">
              {value.toFixed(0)}
            </text>
          </g>
        );
      })}

      {/* Box plots */}
      {transformedData.map((item, idx) => {
        const xPos = padding + 100 + idx * 120;
        const scaleY = (val: number) => height - 40 - ((val - yMin) / yRange) * (height - 60);

        // Whiskers
        const minY = scaleY(item.min);
        const maxY = scaleY(item.max);

        // Box
        const q1Y = scaleY(item.q1);
        const q3Y = scaleY(item.q3);
        const medianY = scaleY(item.median);

        return (
          <g key={idx}>
            {/* Whisker line (min to max) */}
            <line
              x1={xPos}
              y1={minY}
              x2={xPos}
              y2={maxY}
              stroke="#9ca3af"
              strokeWidth={1}
            />

            {/* Min cap */}
            <line x1={xPos - 10} y1={minY} x2={xPos + 10} y2={minY} stroke="#9ca3af" strokeWidth={1} />

            {/* Max cap */}
            <line x1={xPos - 10} y1={maxY} x2={xPos + 10} y2={maxY} stroke="#9ca3af" strokeWidth={1} />

            {/* Box (Q1 to Q3) */}
            <rect
              x={xPos - boxWidth / 2}
              y={Math.min(q1Y, q3Y)}
              width={boxWidth}
              height={Math.abs(q3Y - q1Y)}
              fill="#3b82f6"
              fillOpacity={0.3}
              stroke="#3b82f6"
              strokeWidth={2}
            />

            {/* Median line */}
            <line
              x1={xPos - boxWidth / 2}
              y1={medianY}
              x2={xPos + boxWidth / 2}
              y2={medianY}
              stroke="#ef4444"
              strokeWidth={2}
            />

            {/* X-axis label */}
            <text
              x={xPos}
              y={height - 15}
              textAnchor="middle"
              className="text-xs fill-gray-600"
            >
              {item.label.split('\n')[0].substring(0, 8)}...
            </text>
            <text
              x={xPos}
              y={height - 5}
              textAnchor="middle"
              className="text-xs font-semibold fill-gray-700"
            >
              {item.label.split('\n')[1]}
            </text>
          </g>
        );
      })}
    </svg>
  );
};

export default BoxPlotChart;