interface PercentileResponse {
  [method: string]: Array<{ percentile: string; value: number }>;
}

export default function transformPercentileData(data: PercentileResponse) {
  const percentiles = ["P50", "P75", "P90", "P95", "P99"];

  return percentiles.map(p => {
    const dataPoint: any = { percentile: p };

    Object.entries(data).forEach(([method, points]) => {
      const point = points.find(pt => pt.percentile === p);
      dataPoint[method] = point?.value || 0;
    });

    return dataPoint;
  });
}