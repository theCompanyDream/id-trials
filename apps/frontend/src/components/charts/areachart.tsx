import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

interface AreaChartDataPoint {
  time_bucket: string;
  [key: string]: string | number;
}

interface AreaSeries {
  dataKey: string;
  name: string;
  fill: string;
  stroke: string;
}

interface AreaChartComponentProps {
  data: AreaChartDataPoint[];
  series: AreaSeries[];
  width?: number | string;
  height?: number | string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  type?: 'monotone' | 'linear' | 'natural';
  stackId?: string; // For stacked areas
}

export default function AreaChartComponent({
  data,
  series,
  width = '100%',
  height = 400,
  showGrid = true,
  showLegend = true,
  showTooltip = true,
  type = 'monotone',
  stackId = 'stack'
}: AreaChartComponentProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <AreaChart data={data}>
        {showGrid && <CartesianGrid strokeDasharray="3 3" />}

        <XAxis
          dataKey="time_bucket"
          tickFormatter={(value) => new Date(value).toLocaleTimeString()}
        />

        <YAxis />

        {showTooltip && <Tooltip />}
        {showLegend && <Legend />}

        {series.map((s, index) => (
          <Area
            key={index}
            type={type}
            dataKey={s.dataKey}
            name={s.name}
            fill={s.fill}
            stroke={s.stroke}
            stackId={stackId}
            fillOpacity={0.6}
          />
        ))}
      </AreaChart>
    </ResponsiveContainer>
  );
}