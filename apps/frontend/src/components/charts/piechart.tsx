import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface PieData {
  name: string;
  value: number;
}

interface PieChartComponentProps {
  data: PieData[];
  colors?: string[];
  width?: number | string;
  height?: number | string;
  innerRadius?: number;
  outerRadius?: number;
  showLegend?: boolean;
  showTooltip?: boolean;
  labelLine?: boolean;
  dataKey?: string;
}

const DEFAULT_COLORS = [
  '#0088FE', '#00C49F', '#FFBB28', '#FF8042',
  '#8884D8', '#82CA9D', '#FFC658', '#FF6B9D'
];

export default function PieChartComponent({
  data,
  colors = DEFAULT_COLORS,
  width = '100%',
  height = 400,
  innerRadius = 0,
  outerRadius = 80,
  showLegend = true,
  showTooltip = true,
  labelLine = false,
  dataKey = 'value'
}: PieChartComponentProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={labelLine}
          outerRadius={outerRadius}
          innerRadius={innerRadius}
          fill="#8884d8"
          dataKey={dataKey}
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
          ))}
        </Pie>
        {showTooltip && <Tooltip />}
        {showLegend && <Legend />}
      </PieChart>
    </ResponsiveContainer>
  );
}