import {
  ScatterChart,
  Scatter,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ZAxis
} from 'recharts';

interface BubbleDataPoint {
  x: number;           // DB duration
  y: number;           // Handler duration
  z: number;           // Response size (bubble size)
  name: string;        // ID type or label
  [key: string]: any;  // Additional metadata
}

interface BubbleSeries {
  data: BubbleDataPoint[];
  name: string;
  color: string;
}

interface BubbleChartProps {
  series: BubbleSeries[];
  width?: number | string;
  height?: number | string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  xAxisLabel?: string;
  yAxisLabel?: string;
  title?: string;
  minBubbleSize?: number;
  maxBubbleSize?: number;
}

export default function BubbleChart({
  series,
  width = '100%',
  height = 400,
  showGrid = true,
  showLegend = true,
  showTooltip = true,
  xAxisLabel,
  yAxisLabel,
  title,
  minBubbleSize = 20,
  maxBubbleSize = 400
}: BubbleChartProps) {
  return (
    <div style={{ width }}>
      {title && <h3 style={{ marginBottom: '10px', textAlign: 'center' }}>{title}</h3>}

      <ResponsiveContainer width="100%" height={height}>
        <ScatterChart margin={{ top: 20, right: 30, left: 20, bottom: 20 }}>
          {showGrid && <CartesianGrid strokeDasharray="3 3" />}

          <XAxis
            type="number"
            dataKey="x"
            name={xAxisLabel || 'X'}
            label={xAxisLabel ? { value: xAxisLabel, position: 'insideBottom', offset: -10 } : undefined}
          />

          <YAxis
            type="number"
            dataKey="y"
            name={yAxisLabel || 'Y'}
            label={yAxisLabel ? { value: yAxisLabel, angle: -90, position: 'insideLeft' } : undefined}
          />

          <ZAxis
            type="number"
            dataKey="z"
            range={[minBubbleSize, maxBubbleSize]}
            name="Size"
          />

          {showTooltip && <Tooltip cursor={{ strokeDasharray: '3 3' }} />}
          {showLegend && <Legend />}

          {series.map((s) => (
            <Scatter
              key={s.name}
              name={s.name}
              data={s.data}
              fill={s.color}
            />
          ))}
        </ScatterChart>
      </ResponsiveContainer>
    </div>
  );
}