import { useEffect, useState } from "react";
import { fetchData } from "./api";

function App() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    fetchData().then(setData).catch(console.error);
  }, []);

  return (
    <div>
      <h1>Connexion API Go</h1>
      {data ? <pre>{JSON.stringify(data, null, 2)}</pre> : <p>Chargement...</p>}
    </div>
  );
}

export default App;
