import {
  ResponsiveContainer,
  Tooltip
} from 'recharts';

interface HeatMapDataPoint {
  date: string; // 'YYYY-MM-DD'
  value: number;
}

interface HeatMapCalendarProps {
  data: HeatMapDataPoint[];
  width?: number | string;
  height?: number | string;
  colorScale?: string[];
  cellSize?: number;
  cellGap?: number;
  showTooltip?: boolean;
  title?: string;
}

const DEFAULT_COLOR_SCALE = [
  '#ebedf0', // 0 - empty
  '#9be9a8', // low
  '#40c463', // medium-low
  '#30a14e', // medium-high
  '#216e39'  // high
];

export default function HeatMapCalendar({
  data,
  width = '100%',
  height = 200,
  colorScale = DEFAULT_COLOR_SCALE,
  cellSize = 12,
  cellGap = 3,
  showTooltip = true,
  title
}: HeatMapCalendarProps) {
  // Group data by week and day
  const getWeekData = () => {
    const dataMap = new Map(data.map(d => [d.date, d.value]));
    const weeks: Array<Array<{ date: string; value: number; dayOfWeek: number }>> = [];

    // Get date range
    const dates = data.map(d => new Date(d.date)).sort((a, b) => a.getTime() - b.getTime());
    if (dates.length === 0) return weeks;

    const startDate = new Date(dates[0]);
    const endDate = new Date(dates[dates.length - 1]);

    // Start from Sunday of the first week
    const firstSunday = new Date(startDate);
    firstSunday.setDate(startDate.getDate() - startDate.getDay());

    let currentWeek: Array<{ date: string; value: number; dayOfWeek: number }> = [];
    let currentDate = new Date(firstSunday);

    while (currentDate <= endDate) {
      const dateString = currentDate.toISOString().split('T')[0];
      const value = dataMap.get(dateString) || 0;
      const dayOfWeek = currentDate.getDay();

      currentWeek.push({ date: dateString, value, dayOfWeek });

      if (dayOfWeek === 6) { // Saturday - end of week
        weeks.push(currentWeek);
        currentWeek = [];
      }

      currentDate.setDate(currentDate.getDate() + 1);
    }

    if (currentWeek.length > 0) {
      weeks.push(currentWeek);
    }

    return weeks;
  };

  const getColor = (value: number) => {
    if (value === 0) return colorScale[0];

    const max = Math.max(...data.map(d => d.value));
    const normalized = value / max;

    if (normalized < 0.25) return colorScale[1];
    if (normalized < 0.5) return colorScale[2];
    if (normalized < 0.75) return colorScale[3];
    return colorScale[4];
  };

  const weeks = getWeekData();
  const totalWidth = weeks.length * (cellSize + cellGap);
  const totalHeight = 7 * (cellSize + cellGap);

  return (
    <div style={{ width, position: 'relative' }}>
      {title && <h3 style={{ marginBottom: '10px' }}>{title}</h3>}

      <svg width={totalWidth} height={totalHeight}>
        {weeks.map((week, weekIndex) => (
          <g key={weekIndex} transform={`translate(${weekIndex * (cellSize + cellGap)}, 0)`}>
            {week.map((day) => (
              <g key={day.date}>
                <rect
                  x={0}
                  y={day.dayOfWeek * (cellSize + cellGap)}
                  width={cellSize}
                  height={cellSize}
                  fill={getColor(day.value)}
                  stroke="#ddd"
                  strokeWidth={0.5}
                  rx={2}
                >
                  <title>{`${day.date}: ${day.value}`}</title>
                </rect>
              </g>
            ))}
          </g>
        ))}
      </svg>

      {/* Legend */}
      <div style={{ marginTop: '10px', display: 'flex', alignItems: 'center', gap: '5px', fontSize: '12px' }}>
        <span>Less</span>
        {colorScale.map((color, i) => (
          <div
            key={i}
            style={{
              width: cellSize,
              height: cellSize,
              backgroundColor: color,
              border: '1px solid #ddd',
              borderRadius: '2px'
            }}
          />
        ))}
        <span>More</span>
      </div>
    </div>
  );
}