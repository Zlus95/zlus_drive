import "./index.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Navigation from "./Navigation/Navigation";
import { DialogProvider } from "./providers/DialogProvider";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <DialogProvider>
        <Navigation />
      </DialogProvider>
    </QueryClientProvider>
  );
}

export default App;
