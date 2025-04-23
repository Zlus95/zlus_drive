import React from "react";
import Header from "../../ui-kit/Header";
import api from "../../api";
import { useQuery } from "@tanstack/react-query";

function useGetFiles() {
  return useQuery({
    queryKey: ["filesList"],
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const response = await api.get("/files");
      return response.data;
    },
  });
}

const Home = () => {
  const { isLoading, data, error, isError, isSuccess } = useGetFiles();

  return (
    <div>
      <Header />
    </div>
  );
};

export default Home;
