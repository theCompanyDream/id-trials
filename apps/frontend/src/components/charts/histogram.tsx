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

interface HistogramBin {
  range: string;        // e.g., "0-50", "50-100"
  count: number;
  frequency?: number;   // Optional percentage
}

interface HistogramProps {
  data: HistogramBin[];
  width?: number | string;
  height?: number | string;
  barColor?: string;
  showGrid?: boolean;
  showLegend?: boolean;
  showTooltip?: boolean;
  xAxisLabel?: string;
  yAxisLabel?: string;
  title?: string;
  barSize?: number;
  dataKey?: string;
}

export default function Histogram({
  data,
  width = '100%',
  height = 400,
  barColor = '#8884d8',
  showGrid = true,
  showLegend = false,
  showTooltip = true,
  xAxisLabel,
  yAxisLabel = 'Frequency',
  title,
  barSize = 40,
  dataKey = 'count'
}: HistogramProps) {
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
            dataKey="range"
            label={xAxisLabel ? { value: xAxisLabel, position: 'insideBottom', offset: -10 } : undefined}
          />

          <YAxis
            label={yAxisLabel ? { value: yAxisLabel, angle: -90, position: 'insideLeft' } : undefined}
          />

          {showTooltip && <Tooltip />}
          {showLegend && <Legend />}

          <Bar
            dataKey={dataKey}
            fill={barColor}
            maxBarSize={barSize}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

// Helper function to create bins from raw data
export function createHistogramBins(
  values: number[],
  numBins: number = 10
): HistogramBin[] {
  if (values.length === 0) return [];

  const min = Math.min(...values);
  const max = Math.max(...values);
  const binWidth = (max - min) / numBins;

  const bins: HistogramBin[] = [];

  for (let i = 0; i < numBins; i++) {
    const binStart = min + (i * binWidth);
    const binEnd = binStart + binWidth;

    const count = values.filter(v => {
      if (i === numBins - 1) {
        // Last bin includes the max value
        return v >= binStart && v <= binEnd;
      }
      return v >= binStart && v < binEnd;
    }).length;

    bins.push({
      range: `${Math.round(binStart)}-${Math.round(binEnd)}`,
      count,
      frequency: (count / values.length) * 100
    });
  }

  return bins;
}

// Helper for custom bin ranges
export function createCustomHistogramBins(
  values: number[],
  binRanges: Array<[number, number]> // e.g., [[0, 50], [50, 100], [100, 200]]
): HistogramBin[] {
  return binRanges.map(([start, end]) => {
    const count = values.filter(v => v >= start && v < end).length;

    return {
      range: `${start}-${end}`,
      count,
      frequency: (count / values.length) * 100
    };
  });
}