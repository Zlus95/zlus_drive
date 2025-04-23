import React from "react";
import api from "../../api";
import { useQuery } from "@tanstack/react-query";
import Loader from "../../ui-kit/Loader/Loader";
import ErrorText from "../../ui-kit/ErrorText/ErrorText";

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

  if (isError) {
    return <ErrorText error={error} />;
  }
};

export default Home;
