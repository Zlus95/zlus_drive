import React from "react";
import Header from "../../ui-kit/Header";
import api from "../../api";
import { useQuery } from "@tanstack/react-query";
import Loader from "../../ui-kit/Loader/Loader";

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

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-full">
        <Loader text="Loading..." />
      </div>
    );
  }
};

export default Home;
