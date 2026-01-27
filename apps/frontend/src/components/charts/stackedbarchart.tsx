import { BarChart, Bar, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer, CartesianGrid } from 'recharts';

interface StackConfig {
  dataKey: string;
  name: string;
  fill: string;
}

interface StackedBarChartProps {
  data: any[];
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
  stacks,
  width = '100%',
  height = 400,
  xAxisLabel,
  yAxisLabel = 'auto',
  showGrid = false,
  showLegend = true,
  showTooltip = true,
  layout = 'vertical'
}: StackedBarChartProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <BarChart data={data} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
        {showGrid && <CartesianGrid strokeDasharray="3 3" />}

        <XAxis dataKey={xAxisLabel} />
        <YAxis width={yAxisLabel} />

        {showTooltip && <Tooltip />}
        {showLegend && <Legend />}

        {stacks.map((stack, index) => (
          <Bar
            dataKey={stack.dataKey}
            stackId={`stack`}
            fill={stack.fill}
            name={stack.name}
            background
          />
        ))}

      </BarChart>
    </ResponsiveContainer>
  );
}