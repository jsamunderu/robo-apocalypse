import { useState, useEffect, useCallback } from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from "react-chartjs-2";

ChartJS.register(ArcElement, Tooltip, Legend);

const InfectedChart = (props) => {
  const [stats, setStats] = useState([]);
  const getStats = useCallback(async () => {
    const response = await fetch("http://localhost:8080/survivors/stats");
    const values = await response.json();
    setStats(values);
    console.log(values);
  }, []);

  useEffect(() => {
    getStats();
  }, [getStats]);

  const data = {
    labels: ["Infected", "Healthy"],
    datasets: [
      {
        label: "Percentage infected and healthy",
        data: [stats.infectedPercentage, stats.healthyPercentage],
        backgroundColor: ["rgba(255, 99, 132, 0.2)", "rgba(54, 162, 235, 0.2)"],
        borderColor: ["rgba(255, 99, 132, 1)", "rgba(54, 162, 235, 1)"],
        borderWidth: 1,
      },
    ],
  };

  return <Pie data={data} />;
};

export default InfectedChart;
