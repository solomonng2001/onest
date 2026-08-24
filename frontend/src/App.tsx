import { useState } from "react";
import ServiceTabs, { type Tab } from "./components/ServiceTabs";
import UENPage from "./pages/UENPage";
import WeatherPage from "./pages/WeatherPage";
import "./App.css";

function App() {
  const [activeTab, setActiveTab] = useState<Tab>("weather");

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="container">
          <h1>OneST Web Portal</h1>
        </div>
      </header>

      <main className="container main-content">
        <ServiceTabs
          activeTab={activeTab}
          onTabChange={setActiveTab}
        />

        {activeTab === "weather" ? <WeatherPage /> : <UENPage />}
      </main>
    </div>
  );
}

export default App;
