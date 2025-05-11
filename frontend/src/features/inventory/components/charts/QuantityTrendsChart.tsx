import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

interface QuantityChartData {
  name: string;
  value: number;
  unit?: string;
  id: string;
}

interface QuantityTrendsChartProps {
  data: QuantityChartData[];
}

export function QuantityTrendsChart({ data }: QuantityTrendsChartProps) {
  // If no data or empty data, show a message
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400">
        No quantity data available
      </div>
    );
  }

  // Prepare data for display - truncate long names
  const displayData = data.map((item) => ({
    ...item,
    // Truncate long names for better display
    displayName:
      item.name.length > 15 ? `${item.name.substring(0, 12)}...` : item.name,
  }));

  return (
    <ResponsiveContainer width="100%" height="100%">
      <BarChart
        data={displayData}
        layout="vertical"
        margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid strokeDasharray="3 3" horizontal={false} />
        <XAxis
          type="number"
          domain={[0, "dataMax + 5"]}
          tickFormatter={(value) => Math.floor(value).toString()}
        />
        <YAxis
          dataKey="displayName"
          type="category"
          width={100}
          tick={{ fontSize: 12 }}
        />
        <Tooltip
          formatter={(value: number, name: string, props) => {
            if (props.payload) {
              const item = props.payload;
              const unitText = item.unit ? ` ${item.unit}` : "";
              return [`${value}${unitText}`, item.name];
            }
            return [`${value}`, name];
          }}
          contentStyle={{
            backgroundColor: "white",
            border: "1px solid #ddd",
            borderRadius: "4px",
            padding: "8px",
            fontSize: "12px",
          }}
        />
        <Bar
          dataKey="value"
          fill="#2563eb"
          name="Quantity"
          radius={[0, 4, 4, 0]}
        />
      </BarChart>
    </ResponsiveContainer>
  );
}
