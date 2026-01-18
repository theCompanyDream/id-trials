import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

interface PercentileDataPoint {
  group: string;
  [key: string]: string | number; // p50, p95, p99, etc.
}

interface PercentileConfig {
  key: string;
  name: string;
  color: string;
}

interface PercentileComparisonChartProps {
  data: PercentileDataPoint[];
  percentiles: PercentileConfig[];
  width?: number | string;
  height?: number | string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  xAxisLabel?: string;
  yAxisLabel?: string;
  barSize?: number;
  title?: string;
}

const DEFAULT_PERCENTILES: PercentileConfig[] = [
  { key: 'p50', name: 'P50 (Median)', color: '#8884d8' },
  { key: 'p95', name: 'P95', color: '#82ca9d' },
  { key: 'p99', name: 'P99', color: '#ffc658' }
];

export default function PercentileComparisonChart({
  data,
  percentiles = DEFAULT_PERCENTILES,
  width = '100%',
  height = 400,
  showGrid = true,
  showLegend = true,
  showTooltip = true,
  xAxisLabel,
  yAxisLabel,
  barSize = 20,
  title
}: PercentileComparisonChartProps) {
  return (
    <div style={{ width }}>
      {title && <h3 style={{ marginBottom: '10px', textAlign: 'center' }}>{title}</h3>}

      <ResponsiveContainer width="100%" height={height}>
        <BarChart
          data={data}
          margin={{ top: 20, right: 30, left: 20, bottom: 20 }}
        >
          {showGrid && <CartesianGrid strokeDasharray="3 3" />}

          <XAxis
            dataKey="group"
            label={xAxisLabel ? { value: xAxisLabel, position: 'insideBottom', offset: -10 } : undefined}
          />

          <YAxis
            label={yAxisLabel ? { value: yAxisLabel, angle: -90, position: 'insideLeft' } : undefined}
          />

          {showTooltip && <Tooltip />}
          {showLegend && <Legend />}

          {percentiles.map((percentile) => (
            <Bar
              key={percentile.key}
              dataKey={percentile.key}
              name={percentile.name}
              fill={percentile.color}
              maxBarSize={barSize}
            />
          ))}
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}