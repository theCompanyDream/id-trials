import { BarChart, Bar, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer, CartesianGrid } from 'recharts';

interface StackConfig {
  dataKey: string;
  name: string;
  fill: string;
}

interface StackedBarChartProps {
  data: any[];
  xAxisKey: string;
  stacks: StackConfig[];
  width?: number | string;
  height?: number | string;
  xAxisLabel?: string;
  yAxisLabel?: string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  layout?: 'horizontal' | 'vertical';
}

export default function StackedBarChart({
  data,
  xAxisKey,
  stacks,
  width = '100%',
  height = 400,
  xAxisLabel,
  yAxisLabel,
  showGrid = false,
  showLegend = true,
  showTooltip = true,
  layout = 'vertical'
}: StackedBarChartProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <BarChart data={data} layout={layout}>
        {showGrid && <CartesianGrid strokeDasharray="3 3" />}

        {layout === 'vertical' ? (
          <>
            <XAxis
              dataKey={xAxisKey}
              label={xAxisLabel ? { value: xAxisLabel, position: 'insideBottom', offset: -5 } : undefined}
            />
            <YAxis
              label={yAxisLabel ? { value: yAxisLabel, angle: -90, position: 'insideLeft' } : undefined}
            />
          </>
        ) : (
          <>
            <XAxis
              type="number"
              label={xAxisLabel ? { value: xAxisLabel, position: 'insideBottom', offset: -5 } : undefined}
            />
            <YAxis
              type="category"
              dataKey={xAxisKey}
              label={yAxisLabel ? { value: yAxisLabel, angle: -90, position: 'insideLeft' } : undefined}
            />
          </>
        )}

        {showTooltip && <Tooltip />}
        {showLegend && <Legend />}

        {stacks.map((stack, index) => (
          <Bar
            key={stack.dataKey}
            dataKey={stack.dataKey}
            stackId={`stack:${index}`}
            fill={stack.fill}
            name={stack.name}
          />
        ))}
      </BarChart>
    </ResponsiveContainer>
  );
}