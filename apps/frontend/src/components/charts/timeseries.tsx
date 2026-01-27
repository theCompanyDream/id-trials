import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

interface TimeSeriesDataPoint {
  timestamp: string | number | Date;
  [key: string]: string | number | Date;
}

interface DataSeries {
  dataKey: string;
  color?: string;
  name?: string;
}

interface TimeSeriesChartProps {
  data: TimeSeriesDataPoint[];
  series: DataSeries[];
  width?: number | string;
  height?: number | string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  color: ColorGamut;
  xAxisKey?: string;
  strokeWidth?: number;
  dotSize?: number;
}

const DEFAULT_COLORS = [
  '#8884d8', '#82ca9d', '#ffc658', '#ff7c7c',
  '#8dd1e1', '#a4de6c', '#d0ed57', '#ffa07a'
];

export default function TimeSeriesChart({
  data,
  series,
  width = '100%',
  height = 400,
  showGrid = true,
  showLegend = true,
  showTooltip = true,
  xAxisKey = 'timestamp',
  strokeWidth = 2,
  dotSize = 3
}: TimeSeriesChartProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <LineChart data={data}>
        {showGrid && <CartesianGrid strokeDasharray="3 3" />}

        <XAxis
          dataKey={xAxisKey}
        />

        <YAxis
          width="auto"
        />

        {showTooltip && <Tooltip />}
        {showLegend && <Legend />}

        {series.map((s, index) => (
          <Line
            type="monotone"
            dataKey={s.dataKey}
            stroke={s.color || DEFAULT_COLORS[index % DEFAULT_COLORS.length]}
            strokeWidth={strokeWidth}
            dot={{ r: dotSize }}
            activeDot={{ r: dotSize + 2 }}
          />
        ))}
      </LineChart>
    </ResponsiveContainer>
  );
}