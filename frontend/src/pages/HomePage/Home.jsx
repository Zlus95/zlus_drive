import Header from "../../ui-kit/Header/Header";

const Home = ({ data }) => {
  const { data: files } = data;

  return (
    <>
      <Header />
      <div>
        {files.map((item) => (
          <div key={item.ID}>
            <img
              src={`${process.env.REACT_APP_API_URL}/${item.Path}`}
              alt={item.Name}
            />
          </div>
        ))}
      </div>
    </>
  );
};

export default Home;
