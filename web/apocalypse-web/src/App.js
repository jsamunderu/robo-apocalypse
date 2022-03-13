import Header from "./components/layout/header/header";
import InfectedChart from "./components/ui/infected-chart/infected-chart";
import SurvivorForm from "./components/ui/survivor-form/survivor-form";
import SurvivorsGrid from "./components/ui/survivors-grid/survivors-grid";

function App() {
  return (
    <div className="App">
      <Header />
      <InfectedChart />
      
    </div>
  );
}

export default App;
