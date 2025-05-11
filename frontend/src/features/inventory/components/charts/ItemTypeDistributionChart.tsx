import {
  Cell,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

interface ItemTypeChartData {
  name: string;
  value: number;
  typeId: string;
}

interface ItemTypeDistributionChartProps {
  data: ItemTypeChartData[];
}

export function ItemTypeDistributionChart({
  data,
}: ItemTypeDistributionChartProps) {
  // Color palette for the pie chart (you can customize these colors)
  const COLORS = [
    "#0088FE",
    "#00C49F",
    "#FFBB28",
    "#FF8042",
    "#8884D8",
    "#82CA9D",
    "#A4DE6C",
    "#D0ED57",
    "#FAD000",
    "#F66D44",
  ];

  // If no data or empty data, show a message
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400">
        No item type data available
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <PieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={false}
          outerRadius={80}
          fill="#8884d8"
          dataKey="value"
          label={({ cx, cy, midAngle, outerRadius, percent, name }) => {
            // Only show label if percentage is significant enough (e.g., > 5%)
            if (percent < 0.05) return null;

            const RADIAN = Math.PI / 180;
            const radius = outerRadius + 10;
            const x = cx + radius * Math.cos(-midAngle * RADIAN);
            const y = cy + radius * Math.sin(-midAngle * RADIAN);

            return (
              <text
                x={x}
                y={y}
                fill="#555"
                textAnchor={x > cx ? "start" : "end"}
                dominantBaseline="central"
                fontSize={12}
                fontWeight={500}
              >
                {`${name} (${(percent * 100).toFixed(0)}%)`}
              </text>
            );
          }}
        >
          {data.map((entry, index) => (
            <Cell
              key={`cell-${index}`}
              fill={COLORS[index % COLORS.length]}
              strokeWidth={1}
            />
          ))}
        </Pie>
        <Tooltip
          formatter={(value, name) => [`${value} items`, name]}
          contentStyle={{
            backgroundColor: "white",
            border: "1px solid #ddd",
            borderRadius: "4px",
            padding: "8px",
            fontSize: "12px",
          }}
        />
        <Legend
          layout="horizontal"
          verticalAlign="bottom"
          align="center"
          wrapperStyle={{ padding: "10px 0 0 0" }}
        />
      </PieChart>
    </ResponsiveContainer>
  );
}
