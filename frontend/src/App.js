import "./index.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Navigation from "./Navigation/Navigation";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Navigation />
    </QueryClientProvider>
  );
}

export default App;
