import React, { useState } from "react";
import api from "../../api";
import { useQuery } from "@tanstack/react-query";
import Loader from "../../ui-kit/Loader/Loader";
import ErrorText from "../../ui-kit/ErrorText/ErrorText";
import Home from "./Home";

function useGetFiles(sort) {
  return useQuery({
    queryKey: ["filesList", sort],
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const response = await api.get("/files", {
        params: {
          sort: sort.toString(),
        },
      });
      return response.data;
    },
  });
}

const HomePage = () => {
  const [sort, setSort] = useState(true);
  const { isLoading, data, error, isError, isSuccess } = useGetFiles(sort);

  const toggleSort = () => {
    setSort((prev) => !prev);
  };

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

  if (isSuccess && !isLoading) {
    return <Home data={data} onChangeSort={toggleSort} sort={sort} />;
  }
};

export default HomePage;
