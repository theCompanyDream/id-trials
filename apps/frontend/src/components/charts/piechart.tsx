import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface PieData {
  name: string;
  value: number;
}

interface PieChartComponentProps {
  data: PieData[];
  data2?: PieData[];
  colors?: string[];
  width?: number | string;
  height?: number | string;
  innerRadius?: number;
  outerRadius?: number;
  showLegend?: boolean;
  showTooltip?: boolean;
  labelLine?: boolean;
  dataKey?: string;
  legendLayout?: 'horizontal' | 'vertical';  // ← Add this
  legendAlign?: 'left' | 'center' | 'right'; // ← Add this
}

const DEFAULT_COLORS = [
  '#0088FE', '#00C49F', '#FFBB28', '#FF8042',
  '#8884D8', '#82CA9D', '#FFC658', '#FF6B9D'
];

const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload; // Access full data object

    return (
      <div className="bg-white p-3 border border-gray-300 rounded shadow-lg">
        <p className="font-semibold">{data.name}</p>
        <p className="text-blue-600 font-medium">
          {data.valuePretty || `${data.value} MB`}
        </p>
      </div>
    );
  }
  return null;
};

export default function PieChartComponent({
  data,
  data2,
  colors = DEFAULT_COLORS,
  width = '100%',
  height = 400,
  innerRadius = 20,
  outerRadius = 60,
  showLegend = true,
  showTooltip = true,
  labelLine = false,
  dataKey = 'value',
  legendLayout = 'vertical',      // ← Default to vertical
  legendAlign = 'right'
}: PieChartComponentProps) {
  return (
    <ResponsiveContainer width={width} height={height}>
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={labelLine}
          outerRadius={`${outerRadius}%`}
          innerRadius={`${innerRadius}%`}
          fill="#8884d8"
          dataKey={dataKey}
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
          ))}
        </Pie>
        {data2 && (
          <Pie
            data={data2}
            cx="50%"
            cy="50%"
            labelLine={labelLine}
            innerRadius={`80%`}
            outerRadius={`95%`}
            fill="#82ca9d"
            dataKey={dataKey}
          />
        )}
        {showTooltip && <Tooltip content={<CustomTooltip />} />}
        {showLegend &&
          <Legend
            layout={legendLayout}
            verticalAlign="middle"
            align={legendAlign}
          />}
      </PieChart>
    </ResponsiveContainer>
  );
}