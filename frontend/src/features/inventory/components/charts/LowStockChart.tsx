import {
  Bar,
  BarChart,
  CartesianGrid,
  LabelList,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

interface LowStockChartData {
  name: string;
  value: number;
  unit?: string;
  id: string;
}

interface LowStockChartProps {
  data: LowStockChartData[];
}

export function LowStockChart({ data }: LowStockChartProps) {
  // If no data or empty data, show a message
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400">
        No low stock items
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
          domain={[0, 5]} // Fixed domain for low stock items
          tickCount={6}
          allowDecimals={false}
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
          fill="#ef4444" // Red color for low stock items
          name="Quantity"
          radius={[0, 4, 4, 0]}
        >
          <LabelList
            dataKey="value"
            position="right"
            style={{ fill: "#374151", fontSize: "12px", fontWeight: 500 }}
            formatter={(
              value: number,
              item: LowStockChartData & { displayName: string }
            ) => {
              const unitText = item.unit ? ` ${item.unit}` : "";
              return `${value}${unitText}`;
            }}
          />
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
