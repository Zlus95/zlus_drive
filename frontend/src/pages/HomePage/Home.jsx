import Header from "../../ui-kit/Header/Header";
import Button from "../../ui-kit/Button/Button";

const Home = ({ data }) => {
  const { data: files } = data;

  return (
    <>
      <Header />
      <div className="px-2 space-y-4">
        {files.map((item) => (
          <div
            key={item.ID}
            className="flex items-center gap-4 p-3 border-b border-gray-200 hover:bg-primary"
          >
            <div className="flex-shrink-0">
              <img
                src={`${process.env.REACT_APP_API_URL}/${item.Path}`}
                alt={item.Name}
                className="w-10 h-10 object-cover rounded"
              />
            </div>

            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-white/50 truncate">
                {item.Name}
              </p>
              <p className="text-xs text-gray-500">
                {new Date(item.CreatedAt).toLocaleDateString()} •{" "}
                {(item.Size / 1024).toFixed(1)} KB
              </p>
            </div>

            <div className="flex-shrink-0">
              <Button variant="error">Delete</Button>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default Home;
