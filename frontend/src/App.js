import "./index.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="text-red-500">
        <div>aaaa</div>
      </div>
    </QueryClientProvider>
  );
}

export default App;
