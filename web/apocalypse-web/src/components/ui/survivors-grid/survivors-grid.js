import { useState, useEffect, useCallback } from "react";
import { AgGridReact } from "ag-grid-react";
import "ag-grid-community/dist/styles/ag-grid.css";
import "ag-grid-community/dist/styles/ag-theme-alpine-dark.css";

const SurvivorsGrid = props => {
  const [survivors, setSurvivors] = useState([]);

  const [columnDefs] = useState([
    { field: "name" },
    { field: "age" },
    { field: "gender" },
    { field: "id" },
    { field: "longitude" },
    { field: "latitude" },
    { field: "water" },
    { field: "food" },
    { field: "medication" },
    { field: "ammunition" },
  ]);

  const getSurvivors = useCallback(async () => {
    const response = await fetch("http://localhost:8080/survivors");
    const data = await response.json();
    setSurvivors(data);
    console.log(data);
  }, []);

  useEffect(() => {
    getSurvivors();
  }, [getSurvivors]);
  return (
    <div className="ag-theme-alpine-dark" style={{ height: 400, width: 1800 }}>
      <AgGridReact rowData={survivors} columnDefs={columnDefs} />
    </div>
  );
};

export default SurvivorsGrid;
