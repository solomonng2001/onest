import "./ServiceTabs.css";

export type Tab = "weather" | "uen";

type ServiceTabsProps = {
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
};

function ServiceTabs({
  activeTab,
  onTabChange,
}: ServiceTabsProps) {
  return (
    <nav className="tabs" aria-label="OneST services" role="tablist">
      <button
        id="weather-tab"
        type="button"
        role="tab"
        aria-selected={activeTab === "weather"}
        aria-controls="weather-panel"
        className={
          activeTab === "weather" ? "tab tab-active" : "tab"
        }
        onClick={() => onTabChange("weather")}
      >
        Weather forecast
      </button>

      <button
        id="uen-tab"
        type="button"
        role="tab"
        aria-selected={activeTab === "uen"}
        aria-controls="uen-panel"
        className={activeTab === "uen" ? "tab tab-active" : "tab"}
        onClick={() => onTabChange("uen")}
      >
        UEN validation
      </button>
    </nav>
  );
}

export default ServiceTabs;
